package model

import "time"

type Subscription struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OwnerID       uint      `gorm:"index" json:"owner_id"`
	Brand         string    `gorm:"size:64;index" json:"brand"`
	Title         string    `gorm:"size:128" json:"title"`
	Desc          string    `gorm:"type:text" json:"desc"`
	AccountCipher string    `gorm:"type:text" json:"-"`
	PassCipher    string    `gorm:"type:text" json:"-"`
	PricePerDay   int64     `json:"price_per_day"`
	PricePerWeek  int64     `json:"price_per_week"`
	Status        int8      `gorm:"default:1;index" json:"status"`
	RentedBy      *uint     `json:"rented_by"`
	RentStart     *time.Time `json:"rent_start"`
	RentEnd       *time.Time `json:"rent_end"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// 注：已删除 Owner *User 关联
}

func (Subscription) TableName() string { return "subscriptions" }
