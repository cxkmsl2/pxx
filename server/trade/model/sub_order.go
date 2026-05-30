package model

import "time"

type SubOrder struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderNo        string    `gorm:"size:32;uniqueIndex" json:"order_no"`
	SubscriptionID uint      `gorm:"index" json:"subscription_id"`
	BuyerID        uint      `gorm:"index" json:"buyer_id"`
	Days           int       `json:"days"`
	Amount         int64     `json:"amount"`
	Status         int8      `gorm:"default:1" json:"status"`
	TxID           string    `gorm:"size:64;default:''" json:"tx_id"`
	DisputeReason  string    `gorm:"type:text" json:"dispute_reason"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// 注：已删除 Subscription *Subscription 关联
}

func (SubOrder) TableName() string { return "sub_orders" }
