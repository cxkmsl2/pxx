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

type SubscriptionService struct {
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewSubscriptionService(db *gorm.DB, cm *cache.CacheManager) *SubscriptionService {
	return &SubscriptionService{db: db, cache: cm}
}

func (s *SubscriptionService) List(ctx context.Context, brand string, page, pageSize int) ([]model.Subscription, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var items []model.Subscription
	query := s.db.WithContext(ctx).Model(&model.Subscription{}).Where("status = ?", model.SubStatusAvailable)
	if brand != "" {
		query = query.Where("brand = ?", brand)
	}
	query.Count(&total).Order("created_at DESC").Offset(offset).Limit(limit).Preload("Owner").Find(&items)
	return items, total, nil
}

func (s *SubscriptionService) Create(ctx context.Context, sub *model.Subscription) error {
	// AES encrypt account and password
	accCipher, err := utils.AESEncrypt(sub.AccountCipher)
	if err != nil { return fmt.Errorf("加密失败") }
	passCipher, err := utils.AESEncrypt(sub.PassCipher)
	if err != nil { return fmt.Errorf("加密失败") }
	sub.AccountCipher = accCipher
	sub.PassCipher = passCipher
	sub.Status = model.SubStatusAvailable
	return s.db.WithContext(ctx).Create(sub).Error
}

func (s *SubscriptionService) Rent(ctx context.Context, subID, buyerID uint, days int) (*model.SubOrder, error) {
	// Redis distributed lock
	lockKey := fmt.Sprintf("sub:%d", subID)
	ok, _ := cache.Lock(context.Background(), lockKey, 5*time.Second)
	if !ok {
		return nil, fmt.Errorf("手慢无，已被抢走")
	}
	defer cache.Unlock(context.Background(), lockKey)

	var sub model.Subscription
	if err := s.db.WithContext(ctx).First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("拼卡不存在")
	}
	if sub.Status != model.SubStatusAvailable {
		return nil, fmt.Errorf("该卡已被租用")
	}
	if days <= 0 { days = 1 }

	amount := int64(days) * sub.PricePerDay
	if days >= 7 && sub.PricePerWeek > 0 {
		amount = int64(days/7)*sub.PricePerWeek + int64(days%7)*sub.PricePerDay
	}

	now := time.Now()
	endAt := now.Add(time.Duration(days) * 24 * time.Hour)

	order := &model.SubOrder{
		OrderNo:        fmt.Sprintf("SUB%d%06d", now.UnixMilli()%100000, subID),
		SubscriptionID: subID,
		BuyerID:        buyerID,
		Days:           days,
		Amount:         amount,
		Status:         1,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil { return err }
		return tx.Model(&model.Subscription{}).Where("id = ? AND status = ?", subID, model.SubStatusAvailable).
			Updates(map[string]interface{}{
				"status":    model.SubStatusInUse,
				"rented_by": buyerID,
				"rent_start": now,
				"rent_end":   endAt,
			}).Error
	})
	return order, err
}

func (s *SubscriptionService) GetCredential(ctx context.Context, orderID, userID uint) (account, password string, rentEnd time.Time, err error) {
	var order model.SubOrder
	if err = s.db.WithContext(ctx).Preload("Subscription").First(&order, orderID).Error; err != nil {
		return "", "", time.Time{}, fmt.Errorf("订单不存在")
	}
	if order.BuyerID != userID {
		return "", "", time.Time{}, fmt.Errorf("无权查看")
	}
	if time.Now().After(*order.Subscription.RentEnd) {
		return "", "", time.Time{}, fmt.Errorf("已过期")
	}
	account, err = utils.AESDecrypt(order.Subscription.AccountCipher)
	if err != nil { return "", "", time.Time{}, fmt.Errorf("解密失败") }
	password, err = utils.AESDecrypt(order.Subscription.PassCipher)
	if err != nil { return "", "", time.Time{}, fmt.Errorf("解密失败") }
	return account, password, *order.Subscription.RentEnd, nil
}

func (s *SubscriptionService) Dispute(ctx context.Context, orderID, userID uint, reason string) error {
	return s.db.WithContext(ctx).Model(&model.SubOrder{}).
		Where("id = ? AND buyer_id = ?", orderID, userID).
		Updates(map[string]interface{}{"status": 3, "dispute_reason": reason}).Error
}
