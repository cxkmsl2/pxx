package model

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FromUserID uint      `gorm:"index" json:"from_user_id"`
	ToUserID   uint      `gorm:"index" json:"to_user_id"`
	Content    string    `gorm:"type:text" json:"content"`
	MsgType    string    `gorm:"size:16;default:'text'" json:"msg_type"` // text/image/card
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time `json:"created_at"`

	FromUser *User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   *User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
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
