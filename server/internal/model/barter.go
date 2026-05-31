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

	User *User `gorm:"-" json:"user,omitempty"` // 跨库不可 Preload，由业务层单独查询填充
}

func (BarterItem) TableName() string { return "barter_items" }
