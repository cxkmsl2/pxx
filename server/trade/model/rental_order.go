package model

import "time"

type RentalOrder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RentalID    uint      `gorm:"index" json:"rental_id"`
	RenterID    uint      `gorm:"index" json:"renter_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	TotalAmount int64     `gorm:"default:0" json:"total_amount"`
	Deposit     int64     `gorm:"default:0" json:"deposit"`
	Status      int8      `gorm:"default:1" json:"status"`
	TxID        string    `gorm:"size:64;default:''" json:"tx_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (RentalOrder) TableName() string { return "rental_orders" }
