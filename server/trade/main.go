package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pxx/idl/gen/trade"
	"pxx/pkg/config"
	"pxx/pkg/etcd"
	"pxx/trade/handler"
	tradeModel "pxx/trade/model"
	"pxx/trade/service"

	"pxx/idl/gen/account"
	"pxx/idl/gen/item"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func main() {
	cfg := config.Load()
	port := cfg.ServerPort

	// 1. 连接 MySQL（pxx_trade）
	dsn := cfg.TradeDB.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Trade] db connect failed: %v", err)
	}
	db.AutoMigrate(&tradeModel.Order{}, &tradeModel.GroupBuyOrder{},
		&tradeModel.RentalOrder{}, &tradeModel.BarterOrder{},
		&tradeModel.SubOrder{}, &tradeModel.TccTransaction{})
	log.Println("[Trade] db connected")

	// 2. 从 ETCD 发现 Account 和 Item 服务
	accountConn := mustDial(cfg, "account-service")
	itemConn := mustDial(cfg, "item-service")
	accountCli := account.NewAccountClient(accountConn)
	itemCli := item.NewItemClient(itemConn)

	// 3. 注册到 ETCD
	registry, err := etcd.NewRegistry([]string{cfg.ETCDEndpoints}, etcd.ServiceInfo{
		Name:    "trade-service",
		Address: getLocalAddr(port),
		TTL:     30,
	})
	if err != nil {
		log.Fatalf("[Trade] etcd registry failed: %v", err)
	}
	if err := registry.Register(context.Background()); err != nil {
		log.Fatalf("[Trade] register failed: %v", err)
	}
	defer registry.Deregister(context.Background())

	// 4. 启动 gRPC Server
	svc := service.NewOrderService(db, accountCli, itemCli)
	h := handler.NewTradeHandler(svc)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[Trade] listen failed: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	trade.RegisterTradeServer(grpcServer, h)

	// 5. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[Trade] gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[Trade] serve failed: %v", err)
		}
	}()

	<-quit
	log.Println("[Trade] shutting down...")
	grpcServer.GracefulStop()
	accountConn.Close()
	itemConn.Close()
}

func mustDial(cfg *config.Config, serviceName string) *grpc.ClientConn {
	resolver, err := etcd.NewResolver([]string{cfg.ETCDEndpoints}, serviceName)
	if err != nil {
		log.Fatalf("[Trade] resolve %s failed: %v", serviceName, err)
	}
	addrs, err := resolver.ListEndpoints(context.Background())
	if err != nil || len(addrs) == 0 {
		log.Fatalf("[Trade] no endpoints for %s: %v", serviceName, err)
	}
	resolver.Close()

	conn, err := grpc.Dial(addrs[0],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		log.Fatalf("[Trade] dial %s failed: %v", serviceName, err)
	}
	return conn
}

func getLocalAddr(port string) string {
	host := os.Getenv("HOSTNAME")
	if host == "" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s", host, port)
}
