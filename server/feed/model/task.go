package model

import "time"

const (
	TaskStatusOpen     = 1
	TaskStatusTaken    = 2
	TaskStatusDone     = 3
	TaskStatusCanceled = 4
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PublisherID uint      `gorm:"index" json:"publisher_id"`
	TakerID     *uint     `json:"taker_id"`
	Type        string    `gorm:"size:32;index" json:"type"`
	Title       string    `gorm:"size:128" json:"title"`
	Desc        string    `gorm:"type:text" json:"desc"`
	FromPlace   string    `gorm:"size:128" json:"from_place"`
	ToPlace     string    `gorm:"size:128" json:"to_place"`
	Reward      int64     `json:"reward"`
	Status      int8      `gorm:"default:1;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// 注：已删除 Publisher *User, Taker *User 关联
}

func (Task) TableName() string { return "tasks" }
