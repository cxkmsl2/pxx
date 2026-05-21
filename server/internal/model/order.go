package model

import "time"

type Order struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderNo       string    `gorm:"uniqueIndex;size:32" json:"order_no"`
	BuyerID       uint      `gorm:"index" json:"buyer_id"`
	SellerID      uint      `gorm:"index" json:"seller_id"`
	ProductID     uint      `json:"product_id"`
	Amount        int64     `json:"amount"`
	Status        int8      `gorm:"default:1" json:"status"`
	PayMethod     string    `gorm:"size:16" json:"pay_method"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	DeliveredAt   *time.Time `json:"delivered_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Order) TableName() string { return "orders" }

// 订单状态
const (
	OrderStatusPending   = 1
	OrderStatusPaid      = 2
	OrderStatusDelivered = 3
	OrderStatusCompleted = 4
	OrderStatusCancelled = 5
)
