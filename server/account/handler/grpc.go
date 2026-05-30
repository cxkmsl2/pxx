package handler

import (
	"context"

	"pxx/account/service"
	"pxx/idl/gen/account"
)

type AccountHandler struct {
	account.UnimplementedAccountServer
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) TryFreeze(ctx context.Context, req *account.TryFreezeReq) (*account.TryFreezeResp, error) {
	return h.svc.TryFreeze(ctx, req)
}

func (h *AccountHandler) ConfirmDeduct(ctx context.Context, req *account.ConfirmDeductReq) (*account.ConfirmDeductResp, error) {
	return h.svc.ConfirmDeduct(ctx, req)
}

func (h *AccountHandler) CancelUnfreeze(ctx context.Context, req *account.CancelUnfreezeReq) (*account.CancelUnfreezeResp, error) {
	return h.svc.CancelUnfreeze(ctx, req)
}

func (h *AccountHandler) GetUserInfo(ctx context.Context, req *account.GetUserInfoReq) (*account.GetUserInfoResp, error) {
	return h.svc.GetUserInfo(ctx, req.UserId)
}

func (h *AccountHandler) BatchGetUserInfo(ctx context.Context, req *account.BatchGetUserInfoReq) (*account.BatchGetUserInfoResp, error) {
	users, err := h.svc.BatchGetUserInfo(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	return &account.BatchGetUserInfoResp{Users: users}, nil
}

func (h *AccountHandler) GetBalance(ctx context.Context, req *account.GetBalanceReq) (*account.GetBalanceResp, error) {
	return h.svc.GetBalance(ctx, req.UserId)
}
