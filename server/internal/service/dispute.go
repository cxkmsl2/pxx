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
	// 移除 Preload("Plaintiff") / Preload("Defendant")，因为 User 表在另一个库
	query := s.db.WithContext(ctx).Model(&model.DisputeCase{})
	if status > 0 { query = query.Where("status = ?", status) }
	query.Count(&total).Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&cases)
	return cases, total, nil
}

func (s *DisputeService) Create(ctx context.Context, c *model.DisputeCase) error {
	c.Status = model.CaseStatusReview
	return s.db.WithContext(ctx).Create(c).Error
}

func (s *DisputeService) Vote(ctx context.Context, caseID, voterID uint, voteFor string) (string, error) {
	var res string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var c model.DisputeCase
		// 1. 锁定行
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&c, caseID).Error; err != nil {
			return fmt.Errorf("案件不存在")
		}
		if c.Status != model.CaseStatusPending {
			return fmt.Errorf("案件已结案或不可投票")
		}

		// 2. 检查重复投票
		var existing int64
		tx.Model(&model.CaseVote{}).Where("case_id = ? AND voter_id = ?", caseID, voterID).Count(&existing)
		if existing > 0 {
			return fmt.Errorf("你已经投过票了")
		}

		// 3. 检查投票上限
		if c.PlaintiffVotes + c.DefendantVotes >= 17 {
			return fmt.Errorf("投票已满")
		}

		// 4. 创建投票记录
		vote := &model.CaseVote{CaseID: caseID, VoterID: voterID, VoteFor: voteFor}
		if err := tx.Create(vote).Error; err != nil {
			return err
		}

		// 5. 更新票数并检查胜负
		updates := make(map[string]interface{})
		if voteFor == "plaintiff" {
			c.PlaintiffVotes++
			updates["plaintiff_votes"] = gorm.Expr("plaintiff_votes + 1")
		} else {
			c.DefendantVotes++
			updates["defendant_votes"] = gorm.Expr("defendant_votes + 1")
		}

		if c.PlaintiffVotes >= 9 {
			updates["status"] = model.CaseStatusClosed
			updates["winner"] = "plaintiff"
			res = "winner_plaintiff"
		} else if c.DefendantVotes >= 9 {
			updates["status"] = model.CaseStatusClosed
			updates["winner"] = "defendant"
			res = "winner_defendant"
		} else {
			res = "success"
		}

		if err := tx.Model(&c).Updates(updates).Error; err != nil {
			return err
		}

		// 6. 如果结案，发送系统消息
		if updates["status"] == model.CaseStatusClosed {
			winnerID := c.PlaintiffID
			loserID := c.DefendantID
			if res == "winner_defendant" {
				winnerID = c.DefendantID
				loserID = c.PlaintiffID
			}
			tx.Create(&model.Message{FromUserID: 0, ToUserID: winnerID, Content: "✅ 小法庭判决：您的纠纷（" + c.Title + "）胜诉。", MsgType: "system"})
			tx.Create(&model.Message{FromUserID: 0, ToUserID: loserID, Content: "❌ 小法庭判决：您的纠纷（" + c.Title + "）败诉。", MsgType: "system"})
		}

		return nil
	})

	return res, err
}
