package service

import (
	"context"
	"fmt"

	accountModel "pxx/account/model"
	"pxx/idl/gen/account"
	"pxx/internal/cache"

	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

// generateBranchID 生成分支事务ID
func generateBranchID(txID string, service string) uint64 {
	// 简化实现：使用 hash 生成唯一 branch_id
	var hash uint64 = 5381
	for _, c := range txID + service {
		hash = ((hash << 5) + hash) + uint64(c)
	}
	return hash
}

// TryFreeze 冻结资金（TCC Try 阶段）
// 幂等性保证：利用 tcc_logs 唯一索引防重复执行
// 防悬挂：通过事务保证日志和扣款原子性，若已存在 CANCELLED 记录则 Try 失败
func (s *AccountService) TryFreeze(ctx context.Context, req *account.TryFreezeReq) (*account.TryFreezeResp, error) {
	branchID := generateBranchID(req.TxId, "account")

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 幂等性与防悬挂检查
		var existing accountModel.TCCLog
		if err := tx.Where("tx_id = ? AND branch_id = ?", req.TxId, branchID).First(&existing).Error; err == nil {
			if existing.Status == "CANCELLED" {
				return fmt.Errorf("事务已取消，拒绝 Try")
			}
			return nil // 已执行过，幂等成功
		}

		// 2. 冻结余额（乐观锁：where balance >= amount）
		result := tx.Model(&accountModel.User{}).
			Where("id = ? AND balance >= ? AND status = 1", req.UserId, req.Amount).
			Update("balance", gorm.Expr("balance - ?", req.Amount))
		if result.RowsAffected == 0 {
			return fmt.Errorf("余额不足或用户状态异常")
		}

		// 3. 插入 FROZEN 记录
		log := &accountModel.TCCLog{
			TxID:       req.TxId,
			BranchID:   branchID,
			ActionType: "FREEZE",
			Status:     "FROZEN",
			UserID:     req.UserId,
			Amount:     req.Amount,
			Payload:    "{}",
		}
		return tx.Create(log).Error
	})

	if err != nil {
		return &account.TryFreezeResp{Success: false, Msg: err.Error()}, nil
	}

	return &account.TryFreezeResp{Success: true, BranchId: branchID}, nil
}

// ConfirmDeduct 确认扣款（TCC Confirm 阶段）
func (s *AccountService) ConfirmDeduct(ctx context.Context, req *account.ConfirmDeductReq) (*account.ConfirmDeductResp, error) {
	// 1. 检查日志是否存在
	var log accountModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ?", req.TxId, req.BranchId).First(&log).Error
	if err != nil {
		return &account.ConfirmDeductResp{Success: false, Msg: "事务日志不存在"}, nil
	}

	if log.Status == "CONFIRMED" {
		return &account.ConfirmDeductResp{Success: true, Msg: "幂等返回"}, nil
	}

	// 2. 更新为已确认
	s.db.WithContext(ctx).Model(&log).Update("status", "CONFIRMED")

	// 3. 清除缓存
	cache.GetRedis().Del(ctx, cache.BuildKey("user", log.UserID))

	return &account.ConfirmDeductResp{Success: true}, nil
}

// CancelUnfreeze 解冻资金（TCC Cancel 阶段）
// 防空回滚：若记录不存在则插入 CANCELLED 记录，防悬挂
func (s *AccountService) CancelUnfreeze(ctx context.Context, req *account.CancelUnfreezeReq) (*account.CancelUnfreezeResp, error) {
	branchID := req.BranchId
	if branchID == 0 {
		branchID = generateBranchID(req.TxId, "account")
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var log accountModel.TCCLog
		err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("tx_id = ? AND branch_id = ?", req.TxId, branchID).First(&log).Error
		
		if err != nil {
			// 防空回滚：插入一条 CANCELLED 记录，防止后续延迟的 Try 成功
			emptyLog := &accountModel.TCCLog{
				TxID:       req.TxId,
				BranchID:   branchID,
				ActionType: "UNFREEZE",
				Status:     "CANCELLED",
				Payload:    "{}",
			}
			return tx.Create(emptyLog).Error
		}

		if log.Status == "CANCELLED" {
			return nil // 幂等
		}

		if log.Status == "CONFIRMED" {
			return fmt.Errorf("事务已确认，无法回滚")
		}

		// 只有 FROZEN 状态才需要真实回滚金额
		if log.Status == "FROZEN" {
			tx.Model(&accountModel.User{}).
				Where("id = ?", log.UserID).
				Update("balance", gorm.Expr("balance + ?", log.Amount))
		}

		return tx.Model(&log).Updates(map[string]interface{}{
			"status":      "CANCELLED",
			"action_type": "UNFREEZE",
		}).Error
	})

	if err != nil {
		return &account.CancelUnfreezeResp{Success: false, Msg: err.Error()}, nil
	}

	return &account.CancelUnfreezeResp{Success: true}, nil
}

// GetUserInfo 查询用户信息
func (s *AccountService) GetUserInfo(ctx context.Context, userID uint64) (*account.GetUserInfoResp, error) {
	var u accountModel.User
	err := s.db.WithContext(ctx).First(&u, userID).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &account.GetUserInfoResp{
		Id:         uint64(u.ID),
		Nickname:   u.Nickname,
		Avatar:     u.Avatar,
		Campus:     u.Campus,
		Department: u.Department,
		Credit:     int32(u.Credit),
		Role:       u.Role,
		Balance:    u.Balance,
	}, nil
}

// BatchGetUserInfo 批量查询用户信息
func (s *AccountService) BatchGetUserInfo(ctx context.Context, ids []uint64) ([]*account.GetUserInfoResp, error) {
	var users []accountModel.User
	s.db.WithContext(ctx).Find(&users, ids)
	result := make([]*account.GetUserInfoResp, 0, len(users))
	for _, u := range users {
		result = append(result, &account.GetUserInfoResp{
			Id: uint64(u.ID), Nickname: u.Nickname, Avatar: u.Avatar,
			Campus: u.Campus, Credit: int32(u.Credit), Role: u.Role,
		})
	}
	return result, nil
}

// GetBalance 查询余额
func (s *AccountService) GetBalance(ctx context.Context, userID uint64) (*account.GetBalanceResp, error) {
	var u accountModel.User
	err := s.db.WithContext(ctx).First(&u, userID).Error
	if err != nil {
		return nil, err
	}
	return &account.GetBalanceResp{Balance: u.Balance}, nil
}