package client

import (
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

	c.Account = c.mustDialAccount(cfg)
	c.Item = c.mustDialItem(cfg)
	c.Trade = c.mustDialTrade(cfg)

	log.Println("[Gateway] all gRPC clients ready")
	return c
}

func (c *Clients) mustDialAccount(cfg *config.Config) account.AccountClient {
	conn := c.mustDial(cfg, ServiceAccount)
	return account.NewAccountClient(conn)
}

func (c *Clients) mustDialItem(cfg *config.Config) item.ItemClient {
	conn := c.mustDial(cfg, ServiceItem)
	return item.NewItemClient(conn)
}

func (c *Clients) mustDialTrade(cfg *config.Config) trade.TradeClient {
	conn := c.mustDial(cfg, ServiceTrade)
	return trade.NewTradeClient(conn)
}

func (c *Clients) mustDial(cfg *config.Config, serviceName string) *grpc.ClientConn {
	resolver, err := etcd.NewServiceResolver([]string{cfg.ETCDEndpoints}, serviceName)
	if err != nil {
		log.Fatalf("[Gateway] resolver %s failed: %v", serviceName, err)
	}

	// 等待至少一个端点可用 (缩短启动阻塞)
	var target string
	for i := 0; i < 50; i++ {
		target, _ = resolver.GetEndpoint()
		if target != "" {
			break
		}
		time.Sleep(100 * time.Millisecond) // 每次等100ms，总共最多等5秒
	}

	if target == "" {
		log.Printf("[Gateway] warning: no endpoints for %s yet, will try to dial anyway", serviceName)
		// 如果还没找到，暂时使用服务名作为假目标，补充默认的 8080 端口让 Docker DNS 兜底
		target = fmt.Sprintf("%s:8080", serviceName)
	}
	conn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		log.Printf("[Gateway] dial %s (%s) failed: %v", serviceName, target, err)
	}

	// 注册关闭函数
	c.cleanups = append(c.cleanups, func() {
		if conn != nil { conn.Close() }
		resolver.Close()
	})

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
