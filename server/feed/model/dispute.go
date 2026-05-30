package model

import "time"

const (
	CaseStatusReview  = 0
	CaseStatusPending = 1
	CaseStatusClosed  = 2
)

type DisputeCase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderID        uint      `gorm:"index" json:"order_id"`
	PlaintiffID    uint      `gorm:"index" json:"plaintiff_id"`
	DefendantID    uint      `gorm:"index" json:"defendant_id"`
	Title          string    `gorm:"size:256" json:"title"`
	PlaintiffDesc  string    `gorm:"type:text" json:"plaintiff_desc"`
	DefendantDesc  string    `gorm:"type:text" json:"defendant_desc"`
	PlaintiffVotes int       `gorm:"default:0" json:"plaintiff_votes"`
	DefendantVotes int       `gorm:"default:0" json:"defendant_votes"`
	Status         int8      `gorm:"default:1;index" json:"status"`
	Winner         string    `gorm:"size:16" json:"winner"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// 注：已删除 Plaintiff *User, Defendant *User 关联
}

func (DisputeCase) TableName() string { return "dispute_cases" }

type CaseVote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CaseID    uint      `gorm:"index" json:"case_id"`
	VoterID   uint      `gorm:"index" json:"voter_id"`
	VoteFor   string    `gorm:"size:16" json:"vote_for"`
	CreatedAt time.Time `json:"created_at"`
}

func (CaseVote) TableName() string { return "case_votes" }
