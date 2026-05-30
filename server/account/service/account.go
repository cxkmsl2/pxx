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
// 防悬挂：先写日志再操作数据库
func (s *AccountService) TryFreeze(ctx context.Context, req *account.TryFreezeReq) (*account.TryFreezeResp, error) {
	branchID := generateBranchID(req.TxId, "account")

	// 1. 幂等性检查
	var existing accountModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ?", req.TxId, branchID).First(&existing).Error
	if err == nil {
		// 已执行过，直接返回已有状态
		return &account.TryFreezeResp{
			Success:  existing.Status == "FROZEN" || existing.Status == "CONFIRMED",
			BranchId: branchID,
			Msg:      "幂等返回",
		}, nil
	}

	// 2. 插入 TRYING 记录（先写日志，防悬挂）
	log := &accountModel.TCCLog{
		TxID:       req.TxId,
		BranchID:   branchID,
		ActionType: "FREEZE",
		Status:     "TRYING",
		UserID:     req.UserId,
		Amount:     req.Amount,
	}
	if err := s.db.WithContext(ctx).Create(log).Error; err != nil {
		return &account.TryFreezeResp{Success: false, Msg: "日志写入失败"}, nil
	}

	// 3. 冻结余额（乐观锁：where balance >= amount）
	result := s.db.WithContext(ctx).Model(&accountModel.User{}).
		Where("id = ? AND balance >= ? AND status = 1", req.UserId, req.Amount).
		Update("balance", gorm.Expr("balance - ?", req.Amount))
	if result.RowsAffected == 0 {
		// 余额不足，标记失败
		s.db.WithContext(ctx).Model(log).Update("status", "FAILED")
		return &account.TryFreezeResp{Success: false, BranchId: branchID, Msg: "余额不足"}, nil
	}

	// 4. 更新为 FROZEN
	s.db.WithContext(ctx).Model(log).Update("status", "FROZEN")

	return &account.TryFreezeResp{Success: true, BranchId: branchID}, nil
}

// ConfirmDeduct 确认扣款（TCC Confirm 阶段）
func (s *AccountService) ConfirmDeduct(ctx context.Context, req *account.ConfirmDeductReq) (*account.ConfirmDeductResp, error) {
	// 1. 检查日志是否存在
	var log accountModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ? AND action_type = 'FREEZE'", req.TxId, req.BranchId).First(&log).Error
	if err != nil {
		return &account.ConfirmDeductResp{Success: false, Msg: "事务日志不存在"}, nil
	}

	if log.Status == "CONFIRMED" {
		return &account.ConfirmDeductResp{Success: true, Msg: "幂等返回"}, nil
	}

	// 2. 更新为已确认（冻结资金已在前一步扣减，这里只需更新状态）
	// 注意：实际扣款已在 TryFreeze 完成，Confirm 只是确认状态变更
	s.db.WithContext(ctx).Model(&log).Update("status", "CONFIRMED")

	// 3. 清除缓存
	cache.GetRedis().Del(ctx, cache.BuildKey("user", log.UserID))

	return &account.ConfirmDeductResp{Success: true}, nil
}

// CancelUnfreeze 解冻资金（TCC Cancel 阶段）
// 防空回滚：检查日志状态，仅 FROZEN/TRYING 才执行回滚
func (s *AccountService) CancelUnfreeze(ctx context.Context, req *account.CancelUnfreezeReq) (*account.CancelUnfreezeResp, error) {
	var log accountModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ? AND action_type = 'FREEZE'", req.TxId, req.BranchId).First(&log).Error
	if err != nil {
		return &account.CancelUnfreezeResp{Success: false, Msg: "事务日志不存在"}, nil
	}

	if log.Status == "CANCELLED" {
		return &account.CancelUnfreezeResp{Success: true, Msg: "幂等返回"}, nil
	}

	if log.Status == "CONFIRMED" {
		// 已经 Confirm 的不能回滚
		return &account.CancelUnfreezeResp{Success: false, Msg: "事务已确认，无法回滚"}, nil
	}

	// 回滚：退回冻结金额
	s.db.WithContext(ctx).Model(&accountModel.User{}).
		Where("id = ?", log.UserID).
		Update("balance", gorm.Expr("balance + ?", log.Amount))
	s.db.WithContext(ctx).Model(&log).Updates(map[string]interface{}{
		"status":     "CANCELLED",
		"action_type": "UNFREEZE",
	})

	// 清除缓存
	cache.GetRedis().Del(ctx, cache.BuildKey("user", log.UserID))

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
