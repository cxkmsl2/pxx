package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pxx/internal/cache"
	"pxx/internal/config"
	"pxx/internal/handler"
	"pxx/internal/middleware"
	"pxx/internal/mq"
	"pxx/internal/router"
	"pxx/internal/service"
	"pxx/internal/ws"
	"pxx/pkg/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 1. 连接 MySQL（复用旧的 pxx 数据库，过渡方案）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] connect failed: %v", err)
	}
	db.Exec("SET NAMES utf8mb4")
	log.Println("[DB] connected (legacy pxx db)")

	// 2. Redis
	if err := cache.InitRedis(cfg.RedisAddr, ""); err != nil {
		log.Printf("[Redis] connect failed: %v (continuing)", err)
	} else {
		log.Println("[Redis] connected")
	}
	log.Println("[Cache] ready")

	// 3. MinIO
	utils.InitMinio("pxx-minio:9000", "pxxadmin", "pxxminio123")

	// 4. JWT
	middleware.InitAuth(cfg.JWTSecret)

	// 5. Kafka（如果可用）
	if cfg.KafkaBros != "" {
		mq.InitKafka(cfg.KafkaBros)
		log.Println("[Kafka] ready")
	}

	// 6. WebSocket
	ws.InitHub()

	// 7. 构建服务层
	cm := cache.NewCacheManager()
	services := service.NewServices(db, cm)
	handlers := handler.NewHandlers(services)
	r := router.Setup(handlers)

	// 8. 启动 HTTP
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[Gateway] listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Gateway] listen: %v", err)
		}
	}()

	<-quit
	log.Println("[Gateway] shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
