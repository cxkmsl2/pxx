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

type BarterService struct {
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewBarterService(db *gorm.DB, cm *cache.CacheManager) *BarterService {
	return &BarterService{db: db, cache: cm}
}

func (s *BarterService) List(ctx context.Context, page, pageSize int) ([]model.BarterItem, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var items []model.BarterItem
	s.db.WithContext(ctx).Model(&model.BarterItem{}).Where("status = 1").Count(&total)
	s.db.WithContext(ctx).Where("status = 1").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)
	return items, total, nil
}

func (s *BarterService) Create(ctx context.Context, item *model.BarterItem) error {
	return s.db.WithContext(ctx).Create(item).Error
}

// Agree 双方达成一致，生成置换订单
func (s *BarterService) Agree(ctx context.Context, userBID uint, barterID uint, meetPlace string, meetTime time.Time) (*model.BarterOrder, error) {
	var itemA model.BarterItem
	
	order := &model.BarterOrder{
		OrderNo:   fmt.Sprintf("BR%d%06d", time.Now().UnixMilli()%100000, barterID),
		UserBID:   userBID,
		ItemAID:   barterID,
		ItemBID:   0, // 另一方物品 ID
		Status:    model.BarterStatusAgreed,
		MeetPlace: meetPlace,
		MeetTime:  meetTime,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 检查并锁定物品
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&itemA, barterID).Error; err != nil {
			return fmt.Errorf("置换物品不存在")
		}
		if itemA.Status != 1 {
			return fmt.Errorf("物品已在置换中或已下架")
		}

		// 2. 更新物品状态为“置换中/已达成” (假设 2 表示置换中)
		if err := tx.Model(&itemA).Update("status", 2).Error; err != nil {
			return err
		}

		// 3. 创建订单
		order.UserAID = itemA.UserID
		return tx.Create(order).Error
	})

	return order, err
}