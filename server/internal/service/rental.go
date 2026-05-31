package service

import (
	"context"
	"fmt"
	"time"

	"pxx/internal/cache"
	"pxx/internal/model"
	"pxx/pkg/utils"

	"gorm.io/gorm"
)

type RentalService struct {
	db          *gorm.DB
	cache       *cache.CacheManager
	UserService *UserService
}

func NewRentalService(db *gorm.DB, cm *cache.CacheManager) *RentalService {
	return &RentalService{db: db, cache: cm}
}

func (s *RentalService) List(ctx context.Context, page, pageSize int) ([]model.RentalItem, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var items []model.RentalItem
	s.db.WithContext(ctx).Model(&model.RentalItem{}).Where("status = 1").Count(&total)
	s.db.WithContext(ctx).Where("status = 1").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)
	
	if s.UserService != nil {
		for i := range items {
			if owner, err := s.UserService.GetByID(ctx, items[i].OwnerID); err == nil {
				items[i].Owner = owner
			}
		}
	}
	return items, total, nil
}

func (s *RentalService) Detail(ctx context.Context, id uint) (*model.RentalItem, error) {
	key := cache.BuildKey("rental", id)
	var item model.RentalItem
	err := s.cache.GetOrLoad(ctx, key, &item, 15*time.Minute, func(ctx context.Context) (interface{}, error) {
		var ri model.RentalItem
		if err := s.db.WithContext(ctx).First(&ri, id).Error; err != nil {
			return nil, err
		}
		return &ri, nil
	})
	
	if err == nil && s.UserService != nil && item.OwnerID > 0 {
		if owner, err := s.UserService.GetByID(ctx, item.OwnerID); err == nil {
			item.Owner = owner
		}
	}
	
	return &item, err
}


func (s *RentalService) CreateOrderSimple(ctx context.Context, rentalID, userID uint, days int) (*model.RentalOrder, error) {
	var item model.RentalItem
	now := time.Now()
	order := &model.RentalOrder{
		RentalID:  rentalID,
		RenterID:  userID,
		StartDate: now,
		EndDate:   now.Add(time.Duration(days) * 24 * time.Hour),
		Status:    model.RentalStatusPending,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&item, rentalID).Error; err != nil {
			return fmt.Errorf("租赁物品不存在")
		}
		if item.Status != 1 {
			return fmt.Errorf("物品已被租用或已下架")
		}
		if item.OwnerID == userID {
			return fmt.Errorf("不能租用自己的物品")
		}

		// 更新物品状态为“租用中”
		if err := tx.Model(&item).Update("status", 2).Error; err != nil {
			return err
		}

		order.Deposit = item.Deposit
		order.TotalAmount = int64(days) * item.DailyPrice
		return tx.Create(order).Error
	})

	if err == nil {
		s.cache.Invalidate(ctx, cache.BuildKey("rental", rentalID))
	}
	return order, err
}

func (s *RentalService) CreateOrder(ctx context.Context, order *model.RentalOrder) error {
	return s.db.WithContext(ctx).Create(order).Error
}

func (s *RentalService) Return(ctx context.Context, orderID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.RentalOrder
		if err := tx.First(&order, orderID).Error; err != nil {
			return err
		}
		if order.Status == model.RentalStatusReturned {
			return fmt.Errorf("订单已归还")
		}
		
		// 1. 更新订单状态
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":     model.RentalStatusReturned,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		// 2. 恢复物品状态为“空闲” (status=1)
		return tx.Model(&model.RentalItem{}).Where("id = ?", order.RentalID).Update("status", 1).Error
	})
}