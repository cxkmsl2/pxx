package model

import "time"

type Message struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	FromUserID uint       `gorm:"index" json:"from_user_id"`
	ToUserID   uint       `gorm:"index" json:"to_user_id"`
	Content    string     `gorm:"type:text" json:"content"`
	MsgType    string     `gorm:"size:16;default:'text'" json:"msg_type"`
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time  `json:"created_at"`
	// 注：已删除 FromUser *User, ToUser *User 关联
}

func (Message) TableName() string { return "messages" }

type Conversation struct {
	UserID       uint      `json:"user_id"`
	Nickname     string    `json:"nickname"`
	LastMsg      string    `json:"last_msg"`
	LastTime     time.Time `json:"last_time"`
	Unread       int       `json:"unread"`
	ProductTitle string    `json:"product_title,omitempty"`
}
