package model

import "time"

type GroupBuy struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Title        string    `gorm:"size:128" json:"title"`
	Desc         string    `gorm:"type:text" json:"desc"`
	CoverImage   string    `gorm:"size:512" json:"cover_image"`
	PricePerUnit int64     `json:"price_per_unit"`
	MinPeople    int       `json:"min_people"`
	CurrentPeople int      `gorm:"default:0" json:"current_people"`
	TotalStock   int       `json:"total_stock"`
	SoldCount    int       `gorm:"default:0" json:"sold_count"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	Status       int8      `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (GroupBuy) TableName() string { return "group_buys" }

const (
	GroupBuyStatusPending  = 1
	GroupBuyStatusActive   = 2
	GroupBuyStatusFinished = 3
)
