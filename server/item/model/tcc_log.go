package model

import "time"

// TCCLog 商品库存 TCC 事务日志
type TCCLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TxID        string    `gorm:"size:64;uniqueIndex:idx_tx_branch" json:"tx_id"`
	BranchID    uint64    `gorm:"uniqueIndex:idx_tx_branch" json:"branch_id"`
	ActionType  string    `gorm:"size:16" json:"action_type"`     // LOCK_STOCK / DEDUCT_STOCK / UNLOCK_STOCK
	Status      string    `gorm:"size:16" json:"status"`          // TRYING / LOCKED / CONFIRMED / CANCELLED / FAILED
	ProductID   uint64    `gorm:"default:0" json:"product_id"`
	Payload     string    `gorm:"type:json;default:'{}'" json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (TCCLog) TableName() string { return "tcc_logs" }
