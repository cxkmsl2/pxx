package model

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Type      int8      `gorm:"default:1" json:"type"`
	Title     string    `gorm:"size:128" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Images    string    `gorm:"type:text" json:"images"`
	Tags      string    `gorm:"size:256" json:"tags"`
	ProductID *uint     `json:"product_id"`
	ViewCount int       `gorm:"default:0" json:"view_count"`
	LikeCount int       `gorm:"default:0" json:"like_count"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    *User    `gorm:"-" json:"user,omitempty"`    // 跨库不可 Preload
	Product *Product `gorm:"-" json:"product,omitempty"` // 跨库不可 Preload
}

func (Post) TableName() string { return "posts" }

// 帖子类型
const (
	PostTypeWanted = 1 // 求购悬赏
	PostTypeShare  = 2 // 好物种草
)