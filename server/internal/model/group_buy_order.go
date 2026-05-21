package model

import "time"

type GroupBuyOrder struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	GroupBuyID uint      `gorm:"index" json:"group_buy_id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	IsLeader   bool      `gorm:"default:false" json:"is_leader"`
	Quantity   int       `gorm:"default:1" json:"quantity"`
	Amount     int64     `json:"amount"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (GroupBuyOrder) TableName() string { return "group_buy_orders" }
