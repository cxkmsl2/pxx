package model

import "time"

// TCCLog TCC 事务日志（防悬挂、防空回滚、幂等）
type TCCLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TxID       string    `gorm:"size:64;uniqueIndex:idx_tx_branch" json:"tx_id"`
	BranchID   uint64    `gorm:"uniqueIndex:idx_tx_branch" json:"branch_id"`
	ActionType string    `gorm:"size:16" json:"action_type"`  // FREEZE / DEDUCT / UNFREEZE
	Status     string    `gorm:"size:16" json:"status"`       // TRYING / FROZEN / CONFIRMED / CANCELLED / FAILED
	UserID     uint64    `gorm:"default:0" json:"user_id"`
	Amount     int64     `gorm:"default:0" json:"amount"`
	Payload    string    `gorm:"type:json;default:'{}'" json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TCCLog) TableName() string { return "tcc_logs" }
