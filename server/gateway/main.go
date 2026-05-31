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
	
	"pxx/gateway/client"
	gateway_handler "pxx/gateway/handler"
	pkgConfig "pxx/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[DB] connect failed: %v", err)
		return nil
	}
	db.Exec("SET NAMES utf8mb4")
	return db
}

func main() {
	cfg := config.Load()

	// 1. 连接各微服务数据库 (实现查询侧的数据聚合)
	dbAccount := initDB(cfg.AccountDB.DSN())
	dbItem := initDB(cfg.ItemDB.DSN())
	dbTrade := initDB(cfg.TradeDB.DSN())
	dbFeed := initDB(cfg.FeedDB.DSN())
	log.Println("[Gateway] all databases connected for read-path")

	// 2. 初始化环境
	cache.InitRedis(cfg.RedisAddr, "")
	utils.InitMinio("pxx-minio:9000", "pxxadmin", "pxxminio123")
	middleware.InitAuth(cfg.JWTSecret)
	if cfg.KafkaBros != "" { mq.InitKafka(cfg.KafkaBros) }
	ws.InitHub()

	// 3. 构建服务层 (注入对应的数据源)
	cm := cache.NewCacheManager()
	
	// 注入正确的 DB 实例
	services := service.NewServices(dbAccount, dbItem, dbTrade, dbFeed, cm)
	handlers := handler.NewHandlers(services)

	// 4. 初始化微服务客户端 (增加重试等待逻辑)
	etcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpoints == "" { etcdEndpoints = "etcd:2379" }
	
	pkgCfg := &pkgConfig.Config{ETCDEndpoints: etcdEndpoints}
	
	// 这里会等待其他微服务注册到 ETCD
	gatewayClients := client.InitClients(pkgCfg)
	defer gatewayClients.CloseAll()
	gh := gateway_handler.NewGatewayHandlers(gatewayClients)

	// 5. 核心修复：手动替换 Handler 中的成员，确保指针生效
	handlers.Order.Create = gh.CreateOrder
	handlers.Product.List = gh.ListProducts

	// 6. 启动路由
	r := router.Setup(handlers)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr: addr, Handler: r,
		ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("[Gateway] listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Gateway] listen error: %v", err)
		}
	}()

	<-quit
	log.Println("[Gateway] shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	// 清理数据库连接
	for _, db := range []*gorm.DB{dbAccount, dbItem, dbTrade, dbFeed} {
		if db != nil {
			if sqlDB, err := db.DB(); err == nil {
				sqlDB.Close()
			}
		}
	}
	log.Println("[Gateway] shutdown complete")
}
