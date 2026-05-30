package model

import "time"

type BarterItem struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"index" json:"user_id"`
	Title    string    `gorm:"size:128" json:"title"`
	Desc     string    `gorm:"type:text" json:"desc"`
	Images   string    `gorm:"type:text" json:"images"`
	WantItem string    `gorm:"size:256" json:"want_item"`
	Status   int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 注：已删除 User *User 关联
}

func (BarterItem) TableName() string { return "barter_items" }
