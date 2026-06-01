package main

import (
	"fmt"
	"log"

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

func initDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] connect failed: %v", err)
	}
	db.Exec("SET NAMES utf8mb4")
	return db
}

func main() {
	cfg := config.Load()

	// 1. 初始化各微服务 MySQL
	accountDB := initDB(cfg.AccountDB.DSN())
	itemDB := initDB(cfg.ItemDB.DSN())
	tradeDB := initDB(cfg.TradeDB.DSN())
	feedDB := initDB(cfg.FeedDB.DSN())
	log.Println("[DB] all databases connected")

	// 2. 初始化 Redis
	if err := cache.InitRedis(cfg.RedisAddr, ""); err != nil {
		log.Printf("[Redis] connect failed: %v", err)
	}

	// 3. 初始化缓存和组件
	// cache.InitLocalCache()
	ws.InitHub()
	utils.InitMinio("pxx-minio:9000", "pxxadmin", "pxxminio123")
	cache.InitBloomFilter(1000000, 3)

	cm := cache.NewCacheManager()

	// 4. 初始化 Kafka
	mq.InitKafka(cfg.KafkaBros)
	middleware.InitAuth(cfg.JWTSecret)

	// 5. 组装 Service -> Handler -> Router
	services := service.NewServices(accountDB, itemDB, tradeDB, feedDB, cm)
	handlers := handler.NewHandlers(services)
	r := router.Setup(handlers)

	// 6. 启动
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("[Server] starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[Server] failed: %v", err)
	}
}
