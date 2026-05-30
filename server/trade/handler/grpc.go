package handler

import (
	"context"

	"pxx/idl/gen/trade"
	"pxx/trade/service"
)

type TradeHandler struct {
	trade.UnimplementedTradeServer
	svc *service.OrderService
}

func NewTradeHandler(svc *service.OrderService) *TradeHandler {
	return &TradeHandler{svc: svc}
}

func (h *TradeHandler) CreateOrder(ctx context.Context, req *trade.CreateOrderReq) (*trade.CreateOrderResp, error) {
	return h.svc.CreateOrder(ctx, req)
}

func (h *TradeHandler) GetOrder(ctx context.Context, req *trade.GetOrderReq) (*trade.GetOrderResp, error) {
	return h.svc.GetOrder(ctx, req.OrderId)
}

func (h *TradeHandler) CancelOrder(ctx context.Context, req *trade.CancelOrderReq) (*trade.CancelOrderResp, error) {
	err := h.svc.CancelOrder(ctx, req.OrderId, req.Reason)
	if err != nil {
		return &trade.CancelOrderResp{Success: false, Msg: err.Error()}, nil
	}
	return &trade.CancelOrderResp{Success: true}, nil
}

func (h *TradeHandler) UpdateOrderStatus(ctx context.Context, req *trade.UpdateOrderStatusReq) (*trade.UpdateOrderStatusResp, error) {
	err := h.svc.UpdateOrderStatus(ctx, req.OrderId, req.Status)
	if err != nil {
		return &trade.UpdateOrderStatusResp{Success: false}, nil
	}
	return &trade.UpdateOrderStatusResp{Success: true}, nil
}
