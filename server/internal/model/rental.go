package model

import "time"

type RentalItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OwnerID     uint      `gorm:"index" json:"owner_id"`
	Title       string    `gorm:"size:128" json:"title"`
	Desc        string    `gorm:"type:text" json:"desc"`
	CoverImage  string    `gorm:"size:512" json:"cover_image"`
	Images      string    `gorm:"type:text" json:"images"`
	DailyPrice  int64     `json:"daily_price"`
	WeeklyPrice int64     `json:"weekly_price"`
	Deposit     int64     `json:"deposit"`
	Campus      string    `gorm:"size:64" json:"campus"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Owner *User `gorm:"-" json:"owner,omitempty"` // 跨库不可 Preload，由业务层单独查询填充
}

func (RentalItem) TableName() string { return "rental_items" }
