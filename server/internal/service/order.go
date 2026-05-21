package service

import (
	"context"
	"fmt"
	"time"

	"pxx/internal/model"
	"pxx/internal/mq"

	"pxx/pkg/utils"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) Create(ctx context.Context, buyerID uint, productID uint) (*model.Order, error) {
	var product model.Product
	if err := s.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, fmt.Errorf("商品不存在")
	}
	if product.Status != 1 {
		return nil, fmt.Errorf("商品已下架")
	}

	order := &model.Order{
		OrderNo:   fmt.Sprintf("PX%d%06d", time.Now().UnixMilli()%100000, productID),
		BuyerID:   buyerID,
		SellerID:  product.SellerID,
		ProductID: productID,
		Amount:    product.Price,
		Status:    model.OrderStatusPending,
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}

	// 异步推送到 Kafka
	_ = mq.Publish(ctx, mq.TopicOrderCreated, order)

	// TODO: 15分钟后取消订单（Kafka 延迟消息）
	return order, nil
}

func (s *OrderService) ListMy(ctx context.Context, userID uint, page, pageSize int) ([]model.Order, int64, error) {
	var total int64
	var orders []model.Order
	offset, limit := utils.Paginate(page, pageSize)
	query := s.db.WithContext(ctx).Model(&model.Order{}).
		Where("buyer_id = ? OR seller_id = ?", userID, userID)
	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders)
	return orders, total, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID uint, status int8) error {
	return s.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
