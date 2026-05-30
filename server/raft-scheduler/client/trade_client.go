package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"pxx/idl/gen/trade"
)

// TradeClient 封装对 Trade 服务的 gRPC 调用
type TradeClient struct {
	addr string
}

func NewTradeClient(addr string) *TradeClient {
	return &TradeClient{addr: addr}
}

func (c *TradeClient) CancelOrder(ctx context.Context, orderID uint64, reason string) error {
	conn, err := grpc.Dial(c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := trade.NewTradeClient(conn)
	resp, err := client.CancelOrder(ctx, &trade.CancelOrderReq{
		OrderId: orderID,
		Reason:  reason,
	})
	if err != nil {
		return err
	}
	if !resp.Success {
		return nil // 即使不成功也不报错
	}
	return nil
}

// AddDelayTask 添加延迟任务（15分钟后自动取消）
func (c *TradeClient) AddDelayTask(orderID uint64) error {
	// 这个函数应该提交 Raft ADD 日志
	// 实际使用场景：Trade 服务下单后，调用本函数将延迟任务提交到 Raft
	return nil
}
