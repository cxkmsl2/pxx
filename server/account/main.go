package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"pxx/account/handler"
	accountModel "pxx/account/model"
	"pxx/account/service"
	"pxx/idl/gen/account"
	"pxx/pkg/config"
	"pxx/pkg/etcd"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func main() {
	cfg := config.Load()
	port := cfg.ServerPort

	// 1. 连接 MySQL（pxx_account）
	dsn := cfg.AccountDB.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Account] db connect failed: %v", err)
	}
	db.AutoMigrate(&accountModel.User{}, &accountModel.TCCLog{})
	log.Println("[Account] db connected")

	// 2. 注册到 ETCD
	registry, err := etcd.NewRegistry([]string{cfg.ETCDEndpoints}, etcd.ServiceInfo{
		Name:    "account-service",
		Address: getLocalAddr(port),
		TTL:     30,
	})
	if err != nil {
		log.Fatalf("[Account] etcd registry failed: %v", err)
	}
	if err := registry.Register(context.Background()); err != nil {
		log.Fatalf("[Account] register failed: %v", err)
	}
	defer registry.Deregister(context.Background())

	// 3. 启动 gRPC Server
	svc := service.NewAccountService(db)
	h := handler.NewAccountHandler(svc)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[Account] listen failed: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	account.RegisterAccountServer(grpcServer, h)

	// 4. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[Account] gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[Account] serve failed: %v", err)
		}
	}()

	<-quit
	log.Println("[Account] shutting down...")
	grpcServer.GracefulStop()
}

func getLocalAddr(port string) string {
	host := os.Getenv("HOSTNAME")
	if host == "" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s", host, port)
}
