package model

import "time"

type UserBehavior struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	TargetID   uint      `json:"target_id"`
	TargetType string    `gorm:"size:32" json:"target_type"`
	Action     string    `gorm:"size:16" json:"action"`
	Duration   int       `gorm:"default:0" json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

func (UserBehavior) TableName() string { return "user_behaviors" }
