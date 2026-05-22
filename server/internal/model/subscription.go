package model

import "time"

const (
	SubStatusAvailable = 1 // 空闲中
	SubStatusInUse     = 2 // 使用中
	SubStatusExpired   = 3 // 已过期
)

type Subscription struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OwnerID      uint      `gorm:"index" json:"owner_id"`
	Brand        string    `gorm:"size:64;index" json:"brand"`          // Netflix/B站/百度网盘
	Title        string    `gorm:"size:128" json:"title"`
	Desc         string    `gorm:"type:text" json:"desc"`
	AccountCipher string   `gorm:"type:text" json:"-"`                  // AES 加密的账号
	PassCipher   string    `gorm:"type:text" json:"-"`                  // AES 加密的密码
	PricePerDay  int64     `json:"price_per_day"`                       // 按天价格（分）
	PricePerWeek int64     `json:"price_per_week"`                      // 按周价格（分）
	Status       int8      `gorm:"default:1;index" json:"status"`       // 1空闲 2使用中 3过期
	RentedBy     *uint     `json:"rented_by"`          // 当前租用者ID
	RentStart    *time.Time `json:"rent_start"`
	RentEnd      *time.Time `json:"rent_end"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Owner *User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (Subscription) TableName() string { return "subscriptions" }

type SubOrder struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderNo        string    `gorm:"size:32;uniqueIndex" json:"order_no"`
	SubscriptionID uint      `gorm:"index" json:"subscription_id"`
	BuyerID        uint      `gorm:"index" json:"buyer_id"`
	Days           int       `json:"days"`
	Amount         int64     `json:"amount"`
	Status         int8      `gorm:"default:1" json:"status"`           // 1租用中 2已完成 3已投诉
	DisputeReason  string    `gorm:"type:text" json:"dispute_reason"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Subscription *Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
}

func (SubOrder) TableName() string { return "sub_orders" }
