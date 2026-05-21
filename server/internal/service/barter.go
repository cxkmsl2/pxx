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
	s.db.WithContext(ctx).Where("status = 1").Preload("User").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)
	return items, total, nil
}

func (s *BarterService) Create(ctx context.Context, item *model.BarterItem) error {
	return s.db.WithContext(ctx).Create(item).Error
}

// Agree 双方达成一致，生成置换订单
func (s *BarterService) Agree(ctx context.Context, userBID uint, barterID uint, meetPlace string, meetTime time.Time) (*model.BarterOrder, error) {
	var itemA model.BarterItem
	if err := s.db.WithContext(ctx).First(&itemA, barterID).Error; err != nil {
		return nil, fmt.Errorf("置换物品不存在")
	}

	order := &model.BarterOrder{
		OrderNo:   fmt.Sprintf("BR%d%06d", time.Now().UnixMilli()%100000, barterID),
		UserAID:   itemA.UserID,
		UserBID:   userBID,
		ItemAID:   itemA.ID,
		ItemBID:   0, // 另一方物品 ID
		Status:    model.BarterStatusAgreed,
		MeetPlace: meetPlace,
		MeetTime:  meetTime,
	}
	return order, s.db.WithContext(ctx).Create(order).Error
}
