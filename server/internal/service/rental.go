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
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewRentalService(db *gorm.DB, cm *cache.CacheManager) *RentalService {
	return &RentalService{db: db, cache: cm}
}

func (s *RentalService) List(ctx context.Context, page, pageSize int) ([]model.RentalItem, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var items []model.RentalItem
	s.db.WithContext(ctx).Model(&model.RentalItem{}).Where("status = 1").Count(&total)
	s.db.WithContext(ctx).Where("status = 1").Preload("Owner").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)
	return items, total, nil
}

func (s *RentalService) Detail(ctx context.Context, id uint) (*model.RentalItem, error) {
	key := cache.BuildKey("rental", id)
	var item model.RentalItem
	err := s.cache.GetOrLoad(ctx, key, &item, 15*time.Minute, func(ctx context.Context) (interface{}, error) {
		var ri model.RentalItem
		if err := s.db.WithContext(ctx).Preload("Owner").First(&ri, id).Error; err != nil {
			return nil, err
		}
		return &ri, nil
	})
	return &item, err
}


func (s *RentalService) CreateOrderSimple(ctx context.Context, rentalID, userID uint, days int) (*model.RentalOrder, error) {
	var item model.RentalItem
	if err := s.db.WithContext(ctx).First(&item, rentalID).Error; err != nil {
		return nil, fmt.Errorf("租赁物品不存在")
	}
	now := time.Now()
	order := &model.RentalOrder{
		RentalID:  rentalID,
		RenterID:  userID,
		StartDate: now,
		EndDate:   now.Add(time.Duration(days) * 24 * time.Hour),
		Status:    model.RentalStatusPending,
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (s *RentalService) CreateOrder(ctx context.Context, order *model.RentalOrder) error {
	return s.db.WithContext(ctx).Create(order).Error
}

func (s *RentalService) Return(ctx context.Context, orderID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.RentalOrder{}).Where("id = ?", orderID).Updates(map[string]interface{}{
		"status":     model.RentalStatusReturned,
		"updated_at": now,
	}).Error
}
