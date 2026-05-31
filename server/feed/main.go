package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pxx/feed/consumer"
	feedModel "pxx/feed/model"
	"pxx/feed/service"
	"pxx/pkg/config"
	"pxx/pkg/etcd"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	port := cfg.ServerPort

	// 1. 连接 MySQL（pxx_feed）
	dsn := cfg.FeedDB.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Feed] db connect failed: %v", err)
	}
	db.AutoMigrate(&feedModel.Post{}, &feedModel.Message{}, &feedModel.Task{},
		&feedModel.DisputeCase{}, &feedModel.CaseVote{}, &feedModel.UserBehavior{})
	log.Println("[Feed] db connected")

	// 2. 注册到 ETCD
	registry, err := etcd.NewRegistry([]string{cfg.ETCDEndpoints}, etcd.ServiceInfo{
		Name:    "feed-service",
		Address: getLocalAddr(port),
		TTL:     30,
	})
	if err != nil {
		log.Fatalf("[Feed] etcd registry failed: %v", err)
	}
	if err := registry.Register(context.Background()); err != nil {
		log.Fatalf("[Feed] register failed: %v", err)
	}
	defer registry.Deregister(context.Background())

	// 3. 初始化服务
	postSvc := service.NewPostService(db)

	// 4. 启动 Kafka 消费者
	ctx, cancel := context.WithCancel(context.Background())
	if cfg.KafkaBros != "" {
		consumer.StartFeedConsumer(ctx, cfg.KafkaBros, postSvc)
		log.Println("[Feed] Kafka consumer started")
	}

	log.Printf("[Feed] service started on :%s", port)

	// 5. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Feed] shutting down...")
	cancel()
	registry.Deregister(context.Background())
}

func getLocalAddr(port string) string {
	host := os.Getenv("HOSTNAME")
	if host == "" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s", host, port)
}
