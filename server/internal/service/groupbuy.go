package service

import (
	"context"
	"fmt"
	"time"

	"pxx/internal/cache"
	"pxx/internal/model"
	"pxx/internal/mq"
	"pxx/pkg/utils"

	"gorm.io/gorm"
)

type GroupBuyService struct {
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewGroupBuyService(db *gorm.DB, cm *cache.CacheManager) *GroupBuyService {
	return &GroupBuyService{db: db, cache: cm}
}

func (s *GroupBuyService) List(ctx context.Context, page, pageSize int) ([]model.GroupBuy, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var items []model.GroupBuy
	s.db.WithContext(ctx).Model(&model.GroupBuy{}).Where("status = ?", model.GroupBuyStatusActive).
		Count(&total)
	s.db.WithContext(ctx).Where("status = ?", model.GroupBuyStatusActive).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)
	return items, total, nil
}

func (s *GroupBuyService) Detail(ctx context.Context, id uint) (*model.GroupBuy, error) {
	key := cache.BuildKey("groupbuy", id)
	var gb model.GroupBuy
	err := s.cache.GetOrLoad(ctx, key, &gb, 10*time.Minute, func(ctx context.Context) (interface{}, error) {
		var item model.GroupBuy
		if err := s.db.WithContext(ctx).First(&item, id).Error; err != nil {
			return nil, err
		}
		return &item, nil
	})
	return &gb, err
}

func (s *GroupBuyService) Join(ctx context.Context, groupBuyID, userID uint, quantity int, isLeader bool) error {
	if quantity <= 0 {
		quantity = 1
	}

	var gb model.GroupBuy
	if err := s.db.WithContext(ctx).First(&gb, groupBuyID).Error; err != nil {
		return fmt.Errorf("拼团不存在")
	}
	if gb.Status != model.GroupBuyStatusActive {
		return fmt.Errorf("拼团已结束")
	}
	if gb.SoldCount+quantity > gb.TotalStock {
		return fmt.Errorf("库存不足，仅剩 %d 件", gb.TotalStock-gb.SoldCount)
	}

	order := &model.GroupBuyOrder{
		GroupBuyID: groupBuyID,
		UserID:     userID,
		IsLeader:   isLeader,
		Quantity:   quantity,
		Amount:     int64(quantity) * gb.PricePerUnit,
		Status:     1,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查是否已参与
		var existingCount int64
		tx.Model(&model.GroupBuyOrder{}).
			Where("group_buy_id = ? AND user_id = ?", groupBuyID, userID).Count(&existingCount)
		if existingCount > 0 {
			return fmt.Errorf("你已参与此拼团")
		}

		// 乐观锁/条件更新：检查库存并扣减
		result := tx.Model(&model.GroupBuy{}).
			Where("id = ? AND status = ? AND sold_count + ? <= total_stock", groupBuyID, model.GroupBuyStatusActive, quantity).
			Updates(map[string]interface{}{
				"current_people": gorm.Expr("current_people + 1"),
				"sold_count":     gorm.Expr("sold_count + ?", quantity),
			})

		if result.RowsAffected == 0 {
			return fmt.Errorf("库存不足或拼团已结束")
		}

		return tx.Create(order).Error
	})

	if err == nil {
		s.cache.Invalidate(ctx, cache.BuildKey("groupbuy", groupBuyID))
		_ = mq.Publish(ctx, mq.TopicGroupBuyFlash, order)
	}

	return err
}
