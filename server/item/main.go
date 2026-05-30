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

	"pxx/idl/gen/item"
	itemModel "pxx/item/model"
	"pxx/item/handler"
	"pxx/item/service"
	"pxx/pkg/config"
	"pxx/pkg/etcd"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	cfg := config.Load()
	port := cfg.ServerPort

	// 1. 连接 MySQL（pxx_item）
	dsn := cfg.ItemDB.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Item] db connect failed: %v", err)
	}
	db.AutoMigrate(&itemModel.Product{}, &itemModel.RentalItem{}, &itemModel.BarterItem{},
		&itemModel.Subscription{}, &itemModel.GroupBuy{}, &itemModel.TCCLog{})
	log.Println("[Item] db connected")

	// 2. 注册到 ETCD
	registry, err := etcd.NewRegistry([]string{cfg.ETCDEndpoints}, etcd.ServiceInfo{
		Name:    "item-service",
		Address: getLocalAddr(port),
		TTL:     30,
	})
	if err != nil {
		log.Fatalf("[Item] etcd registry failed: %v", err)
	}
	if err := registry.Register(context.Background()); err != nil {
		log.Fatalf("[Item] register failed: %v", err)
	}
	defer registry.Deregister(context.Background())

	// 3. 启动 gRPC Server
	svc := service.NewItemService(db)
	h := handler.NewItemHandler(svc)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[Item] listen failed: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	item.RegisterItemServer(grpcServer, h)

	// 4. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[Item] gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[Item] serve failed: %v", err)
		}
	}()

	<-quit
	log.Println("[Item] shutting down...")
	grpcServer.GracefulStop()
}

func getLocalAddr(port string) string {
	host := os.Getenv("HOSTNAME")
	if host == "" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s", host, port)
}
