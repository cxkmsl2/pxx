package service

import (
	"context"
	"fmt"

	"pxx/internal/model"

	"gorm.io/gorm"
)

type DisputeService struct {
	db *gorm.DB
}

func NewDisputeService(db *gorm.DB) *DisputeService {
	return &DisputeService{db: db}
}

func (s *DisputeService) Accept(ctx context.Context, caseID uint) error {
	return s.db.WithContext(ctx).Model(&model.DisputeCase{}).Where("id = ? AND status = ?", caseID, model.CaseStatusReview).Update("status", model.CaseStatusPending).Error
}

func (s *DisputeService) ListCases(ctx context.Context, status int8, page, pageSize int) ([]model.DisputeCase, int64, error) {
	var total int64
	var cases []model.DisputeCase
	query := s.db.WithContext(ctx).Model(&model.DisputeCase{}).Preload("Plaintiff").Preload("Defendant")
	if status > 0 { query = query.Where("status = ?", status) }
	query.Count(&total).Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&cases)
	return cases, total, nil
}

func (s *DisputeService) Create(ctx context.Context, c *model.DisputeCase) error {
	c.Status = model.CaseStatusReview
	return s.db.WithContext(ctx).Create(c).Error
}

func (s *DisputeService) Vote(ctx context.Context, caseID, voterID uint, voteFor string) (string, error) {
	var c model.DisputeCase
	if err := s.db.WithContext(ctx).First(&c, caseID).Error; err != nil {
		return "", fmt.Errorf("案件不存在")
	}
	if c.Status != model.CaseStatusPending {
		return "", fmt.Errorf("案件已结案")
	}
	// Check if already voted
	var existing int64
	s.db.WithContext(ctx).Model(&model.CaseVote{}).Where("case_id = ? AND voter_id = ?", caseID, voterID).Count(&existing)
	if existing > 0 {
		return "", fmt.Errorf("你已经投过票了")
	}
	// Check total votes
	totalVotes := c.PlaintiffVotes + c.DefendantVotes
	if totalVotes >= 17 {
		return "", fmt.Errorf("投票已满")
	}
	// Create vote
	vote := &model.CaseVote{CaseID: caseID, VoterID: voterID, VoteFor: voteFor}
	if err := s.db.WithContext(ctx).Create(vote).Error; err != nil {
		return "", err
	}
	// Update count
	if voteFor == "plaintiff" {
		s.db.Model(&c).UpdateColumn("plaintiff_votes", gorm.Expr("plaintiff_votes + 1"))
		c.PlaintiffVotes++
	} else {
		s.db.Model(&c).UpdateColumn("defendant_votes", gorm.Expr("defendant_votes + 1"))
		c.DefendantVotes++
	}
	// Check winner
	if c.PlaintiffVotes >= 9 {
		s.db.Model(&c).Updates(map[string]interface{}{"status": model.CaseStatusClosed, "winner": "plaintiff"})
		// Notify winner
		s.db.WithContext(ctx).Create(&model.Message{FromUserID: 0, ToUserID: c.PlaintiffID, Content: "✅ 小法庭判决：您的纠纷（" + c.Title + "）胜诉，信用分不变。", MsgType: "system"})
		s.db.WithContext(ctx).Create(&model.Message{FromUserID: 0, ToUserID: c.DefendantID, Content: "❌ 小法庭判决：您的纠纷（" + c.Title + "）败诉，信用分 -5。", MsgType: "system"})
		return "winner_plaintiff", nil
	}
	if c.DefendantVotes >= 9 {
		s.db.Model(&c).Updates(map[string]interface{}{"status": model.CaseStatusClosed, "winner": "defendant"})
		// Notify winner
		s.db.WithContext(ctx).Create(&model.Message{FromUserID: 0, ToUserID: c.DefendantID, Content: "✅ 小法庭判决：您的纠纷（" + c.Title + "）胜诉，信用分不变。", MsgType: "system"})
		s.db.WithContext(ctx).Create(&model.Message{FromUserID: 0, ToUserID: c.PlaintiffID, Content: "❌ 小法庭判决：您的纠纷（" + c.Title + "）败诉，信用分 -5。", MsgType: "system"})
		return "winner_defendant", nil
	}
	return "success", nil
}
