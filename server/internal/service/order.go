package service

import (
	"context"
	"fmt"

	"pxx/internal/model"
	tradeModel "pxx/trade/model"

	"pxx/pkg/utils"
	"gorm.io/gorm"
)

type OrderService struct {
	db             *gorm.DB
	ProductService *ProductService
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) Create(ctx context.Context, buyerID uint, productID uint) (*model.Order, error) {
	return nil, fmt.Errorf("Please use Gateway CreateOrder")
}

func (s *OrderService) ListMy(ctx context.Context, userID uint, page, pageSize int) ([]tradeModel.Order, int64, error) {
	var total int64
	var orders []tradeModel.Order
	offset, limit := utils.Paginate(page, pageSize)
	
	query := s.db.WithContext(ctx).Model(&tradeModel.Order{}).
		Where("buyer_id = ? OR seller_id = ?", userID, userID)
	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders)
	return orders, total, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID uint, status int8) error {
	return s.db.WithContext(ctx).Model(&tradeModel.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
