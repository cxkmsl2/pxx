package model

import "time"

const (
	TaskStatusOpen     = 1 // 待接单
	TaskStatusTaken    = 2 // 已接单
	TaskStatusDone     = 3 // 已完成
	TaskStatusCanceled = 4 // 已取消
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PublisherID uint      `gorm:"index" json:"publisher_id"`
	TakerID     *uint     `json:"taker_id"`
	Type        string    `gorm:"size:32;index" json:"type"`        // food/express/print/book
	Title       string    `gorm:"size:128" json:"title"`
	Desc        string    `gorm:"type:text" json:"desc"`
	FromPlace   string    `gorm:"size:128" json:"from_place"`
	ToPlace     string    `gorm:"size:128" json:"to_place"`
	Reward      int64     `json:"reward"`                            // 悬赏金（分）
	Status      int8      `gorm:"default:1;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Publisher *User `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	Taker     *User  `gorm:"foreignKey:TakerID" json:"taker,omitempty"`
}

func (Task) TableName() string { return "tasks" }
