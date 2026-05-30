package handler

import (
	"context"

	"pxx/idl/gen/item"
	"pxx/item/service"
)

type ItemHandler struct {
	item.UnimplementedItemServer
	svc *service.ItemService
}

func NewItemHandler(svc *service.ItemService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

func (h *ItemHandler) TryLockStock(ctx context.Context, req *item.TryLockStockReq) (*item.TryLockStockResp, error) {
	return h.svc.TryLockStock(ctx, req)
}

func (h *ItemHandler) ConfirmDeductStock(ctx context.Context, req *item.ConfirmDeductStockReq) (*item.ConfirmDeductStockResp, error) {
	return h.svc.ConfirmDeductStock(ctx, req)
}

func (h *ItemHandler) CancelUnlockStock(ctx context.Context, req *item.CancelUnlockStockReq) (*item.CancelUnlockStockResp, error) {
	return h.svc.CancelUnlockStock(ctx, req)
}

func (h *ItemHandler) GetProductDetail(ctx context.Context, req *item.GetProductDetailReq) (*item.GetProductDetailResp, error) {
	return h.svc.GetProductDetail(ctx, req.ProductId)
}

func (h *ItemHandler) ListProducts(ctx context.Context, req *item.ListProductsReq) (*item.ListProductsResp, error) {
	return h.svc.ListProducts(ctx, req)
}
