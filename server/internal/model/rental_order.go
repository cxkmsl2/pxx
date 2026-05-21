package model

import "time"

type RentalOrder struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	RentalID    uint       `gorm:"index" json:"rental_id"`
	RenterID    uint       `gorm:"index" json:"renter_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	TotalAmount int64      `json:"total_amount"`
	Deposit     int64      `json:"deposit"`
	Status      int8       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (RentalOrder) TableName() string { return "rental_orders" }

const (
	RentalStatusPending  = 1
	RentalStatusActive   = 2
	RentalStatusReturned = 3
	RentalStatusOverdue  = 4
)
