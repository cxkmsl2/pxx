package model

import "time"

// TccTransaction TCC 全局事务状态表
type TccTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TxID      string    `gorm:"size:64;uniqueIndex" json:"tx_id"`
	Status    string    `gorm:"size:16;default:TRYING" json:"status"` // TRYING / CONFIRMED / CANCELLED
	CreatedAt time.Time `json:"created_at"`
}

func (TccTransaction) TableName() string { return "tcc_transactions" }
