package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"pxx/idl/gen/account"
	"pxx/idl/gen/item"
	"pxx/idl/gen/trade"
	"pxx/pkg/config"
	"pxx/pkg/etcd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Clients 所有下游微服务 gRPC 客户端
type Clients struct {
	Account account.AccountClient
	Item    item.ItemClient
	Trade   trade.TradeClient

	// 关闭函数列表
	cleanups []func()
}

// InitClients 从 ETCD 发现并连接所有微服务
func InitClients(cfg *config.Config) *Clients {
	c := &Clients{}

	c.Account = mustDialAccount(cfg)
	c.Item = mustDialItem(cfg)
	c.Trade = mustDialTrade(cfg)

	log.Println("[Gateway] all gRPC clients ready")
	return c
}

func mustDialAccount(cfg *config.Config) account.AccountClient {
	conn := mustDial(cfg, ServiceAccount)
	return account.NewAccountClient(conn)
}

func mustDialItem(cfg *config.Config) item.ItemClient {
	conn := mustDial(cfg, ServiceItem)
	return item.NewItemClient(conn)
}

func mustDialTrade(cfg *config.Config) trade.TradeClient {
	conn := mustDial(cfg, ServiceTrade)
	return trade.NewTradeClient(conn)
}

func mustDial(cfg *config.Config, serviceName string) *grpc.ClientConn {
	resolver, err := etcd.NewResolver([]string{cfg.ETCDEndpoints}, serviceName)
	if err != nil {
		log.Fatalf("[Gateway] resolve %s failed: %v", serviceName, err)
	}
	addrs, err := resolver.ListEndpoints(context.Background())
	if err != nil || len(addrs) == 0 {
		log.Fatalf("[Gateway] no endpoints for %s: %v", serviceName, err)
	}
	resolver.Close()

	target := addrs[0]
	conn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		log.Fatalf("[Gateway] dial %s (%s) failed: %v", serviceName, target, err)
	}
	return conn
}

// CloseAll 关闭所有连接
func (c *Clients) CloseAll() {
	for _, f := range c.cleanups {
		f()
	}
}

// Helper 类型断言访问器
func (c *Clients) AccountClient() account.AccountClient {
	return c.Account
}
func (c *Clients) ItemClient() item.ItemClient {
	return c.Item
}
func (c *Clients) TradeClient() trade.TradeClient {
	return c.Trade
}

// 注册服务名常量
const (
	ServiceAccount = "account-service"
	ServiceItem    = "item-service"
	ServiceTrade   = "trade-service"
	ServiceFeed    = "feed-service"
)

var serviceNames = []string{"", ServiceAccount, ServiceItem, ServiceTrade, ServiceFeed}

func ServiceAddr(name string, port int) string {
	return fmt.Sprintf("%s:%d", name, port)
}
