package model

import "time"

type BarterOrder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderNo     string    `gorm:"uniqueIndex;size:32" json:"order_no"`
	UserAID     uint      `json:"user_a_id"`
	UserBID     uint      `json:"user_b_id"`
	ItemAID     uint      `json:"item_a_id"`
	ItemBID     uint      `json:"item_b_id"`
	Status      int8      `gorm:"default:1" json:"status"`
	MeetPlace   string    `gorm:"size:128" json:"meet_place"`
	MeetTime    time.Time `json:"meet_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (BarterOrder) TableName() string { return "barter_orders" }

const (
	BarterStatusNegotiating = 1
	BarterStatusAgreed      = 2
	BarterStatusCompleted   = 3
	BarterStatusCancelled   = 4
)
