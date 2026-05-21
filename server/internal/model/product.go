package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SellerID    uint      `gorm:"index" json:"seller_id"`
	Title       string    `gorm:"size:128" json:"title"`
	Desc        string    `gorm:"type:text" json:"desc"`
	Category    string    `gorm:"size:32;index" json:"category"`
	Tag         string    `gorm:"size:32" json:"tag"`
	Price       int64     `json:"price"`
	OriginalPrice int64   `json:"original_price"`
	Images      string    `gorm:"type:text" json:"images"`
	Condition   int8      `gorm:"default:1" json:"condition"`
	Campus      string    `gorm:"size:64" json:"campus"`
	Dormitory   string    `gorm:"size:64" json:"dormitory"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	LikeCount   int       `gorm:"default:0" json:"like_count"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Seller *User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}

func (Product) TableName() string { return "products" }
