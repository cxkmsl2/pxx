package model

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OpenID      string    `gorm:"uniqueIndex;size:64" json:"open_id"`
	Nickname    string    `gorm:"size:64" json:"nickname"`
	Avatar      string    `gorm:"size:512" json:"avatar"`
	Phone       string    `gorm:"size:20" json:"phone"`
	Campus      string    `gorm:"size:64" json:"campus"`
	Department  string    `gorm:"size:64" json:"department"`
	Dormitory   string    `gorm:"size:64" json:"dormitory"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	Balance     int64     `gorm:"default:0" json:"balance"`
	Role        string    `gorm:"size:20;default:user" json:"role"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
