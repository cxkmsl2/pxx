package main

import (
	"fmt"
	"log"
	"context"
	"time"

	"pxx/internal/cache"
	"pxx/internal/config"
	"pxx/internal/model"
	"pxx/pkg/utils"
	"pxx/internal/handler"
	"pxx/internal/middleware"
	"pxx/internal/mq"
	"pxx/internal/router"
	"pxx/internal/ws"
	"pxx/internal/service"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func seedDatabase(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Order{},
		&model.Post{},
		&model.GroupBuy{},
		&model.GroupBuyOrder{},
		&model.RentalItem{},
		&model.RentalOrder{},
		&model.BarterItem{},
		&model.BarterOrder{},
		&model.UserBehavior{},
		&model.Subscription{},
		&model.SubOrder{},
		&model.Task{},
		&model.Message{},
		&model.DisputeCase{},
		&model.CaseVote{},
	)
	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count > 0 { log.Println("[Seed] skip, data exists"); return }
	log.Println("[Seed] inserting...")
	u1 := model.User{OpenID: "test_001", Nickname: "校园小明", Campus: "南校区", Department: "计算机学院"}
	u2 := model.User{OpenID: "test_002", Nickname: "图书馆小张", Campus: "北校区", Department: "软件学院"}
	u3 := model.User{OpenID: "admin_001", Nickname: "校园评审员", Campus: "管理", Department: "信息中心", Role: "admin"}
	db.Create(&u1); db.Create(&u2)
	db.Create(&u3)
	db.Create(&model.Product{SellerID: u1.ID, Title: "考研408全套资料", Desc: "数据结构+组成原理+操作系统+计算机网络，几乎全新", Category: "考研资料", Tag: "计算机", Price: 4500, OriginalPrice: 12000, Campus: "南校区"})
	db.Create(&model.Product{SellerID: u2.ID, Title: "罗技G502鼠标", Desc: "用了半年，功能完好，箱说全", Category: "数码外设", Tag: "鼠标", Price: 15000, OriginalPrice: 39900, Campus: "北校区"})
	db.Create(&model.Product{SellerID: u1.ID, Title: "宿舍用小台灯", Desc: "LED台灯，三档调光", Category: "日用百货", Tag: "台灯", Price: 2500, OriginalPrice: 5000, Campus: "南校区"})
	db.Create(&model.Post{UserID: u1.ID, Type: 1, Title: "求购二手显示器一台", Content: "24寸即可，主要用于宿舍写代码", Tags: "数码,显示器"})
	db.Create(&model.Post{UserID: u2.ID, Type: 2, Title: "好物推荐：图书馆自习必备降噪耳塞", Content: "推荐这款3M耳塞，图书馆自习神器", Tags: "自习,图书馆"})

	// Group buys
	t1 := time.Now()
	db.Create(&model.GroupBuy{Title: "维达抽纸 3层*10包 宿舍拼团", Desc: "柔软亲肤，不掉纸屑，宿舍必备", PricePerUnit: 1500, MinPeople: 3, CurrentPeople: 2, TotalStock: 50, SoldCount: 2, StartAt: t1, EndAt: t1.Add(72 * time.Hour), Status: 2})
	db.Create(&model.GroupBuy{Title: "得力垃圾袋 50只装 阶梯拼团", Desc: "加厚款，承重力强，满5人再减2元", PricePerUnit: 800, MinPeople: 5, CurrentPeople: 3, TotalStock: 100, SoldCount: 3, StartAt: t1, EndAt: t1.Add(48 * time.Hour), Status: 2})

	// Rental items
	db.Create(&model.RentalItem{OwnerID: u1.ID, Title: "佳能 EOS 200D 单反相机", Desc: "入门级单反，配18-55mm镜头，适合社团活动拍摄", DailyPrice: 3500, WeeklyPrice: 20000, Deposit: 50000, Campus: "南校区"})
	db.Create(&model.RentalItem{OwnerID: u2.ID, Title: "任天堂 Switch 游戏机", Desc: "红蓝版，带塞尔达+动森卡带，可租一周送两天", DailyPrice: 2000, WeeklyPrice: 10000, Deposit: 30000, Campus: "北校区"})
	db.Create(&model.RentalItem{OwnerID: u1.ID, Title: "男生正装套装（M码）", Desc: "面试/答辩专用，含西装外套+西裤+领带，九成新", DailyPrice: 1500, WeeklyPrice: 8000, Deposit: 20000, Campus: "南校区"})

	// Barter items
	db.Create(&model.BarterItem{UserID: u1.ID, Title: "考研英语红宝书（九成新）", Desc: "背了前10个list，后面全新", WantItem: "考研数学复习全书"})
		db.Create(&model.BarterItem{UserID: u2.ID, Title: "机械键盘 IKBC C87 青轴", Desc: "用了三个月，换了静电容所以闲置", WantItem: "2K显示器或降噪耳机"})

	// Subscription seed
	sub1 := model.Subscription{OwnerID: u1.ID, Brand: "B站", Title: "B站大会员月卡 独享", Desc: "稳定不掉，可长期租", AccountCipher: "bilibili_user", PassCipher: "b123456", PricePerDay: 100, PricePerWeek: 500, Status: model.SubStatusAvailable}
	sub2 := model.Subscription{OwnerID: u2.ID, Brand: "Netflix", Title: "Netflix 4K 高级套餐", Desc: "4人车还有一个位置，按周租", AccountCipher: "netflix_user", PassCipher: "n123456", PricePerDay: 0, PricePerWeek: 800, Status: model.SubStatusAvailable}
	sub3 := model.Subscription{OwnerID: u1.ID, Brand: "百度网盘", Title: "百度网盘 SVIP 月卡", Desc: "极速下载，考研资料必备", AccountCipher: "baidu_user", PassCipher: "bai123456", PricePerDay: 150, PricePerWeek: 700, Status: model.SubStatusAvailable}
	sub4 := model.Subscription{OwnerID: u2.ID, Brand: "ChatGPT", Title: "ChatGPT Plus 共享", Desc: "GPT-4模型，适合写论文", AccountCipher: "gpt_user", PassCipher: "g123456", PricePerDay: 200, PricePerWeek: 1000, Status: model.SubStatusAvailable}
	// Must encrypt in Create
	for _, s := range []*model.Subscription{&sub1, &sub2, &sub3, &sub4} {
		s.AccountCipher, _ = utils.AESEncrypt(s.AccountCipher)
		s.PassCipher, _ = utils.AESEncrypt(s.PassCipher)
		db.Create(s)
	}

	// Tasks
	db.Create(&model.Task{PublisherID: u1.ID, Type: "food", Title: "帮带一份黄焖鸡米饭", Desc: "微辣，二食堂二楼", FromPlace: "二食堂", ToPlace: "12栋301", Reward: 300})
	db.Create(&model.Task{PublisherID: u2.ID, Type: "express", Title: "菜鸟驿站取快递", Desc: "中通快递，取件码已发", FromPlace: "菜鸟驿站", ToPlace: "图书馆3楼", Reward: 200})
	db.Create(&model.Task{PublisherID: u1.ID, Type: "print", Title: "打印一份实验报告", Desc: "10页黑白双面，装订", FromPlace: "打印店", ToPlace: "教学楼B201", Reward: 100})
	db.Create(&model.Task{PublisherID: u2.ID, Type: "book", Title: "帮还图书馆借书", Desc: "两本，在服务台还", FromPlace: "图书馆", ToPlace: "图书馆", Reward: 200})

	log.Println("[Seed] done")
}

func main() {
	cfg := config.Load()

	// 1. 初始化 MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] connect failed: %v", err)
	}
	db.Exec("SET NAMES utf8mb4")
	log.Println("[DB] connected")
	seedDatabase(db)

	// 2. 初始化 Redis
	if err := cache.InitRedis(cfg.RedisAddr, ""); err != nil {
		log.Printf("[Redis] connect failed: %v (continuing without redis)", err)
	} else {
		log.Println("[Redis] connected")
	}

	// 3. 初始化本地缓存 L2
	cache.InitLocalCache()
	log.Println("[Cache] local cache ready")

	// 4. 初始化布隆过滤器 L1
	// WebSocket Hub
	ws.InitHub()

	// Init MinIO
	utils.InitMinio("pxx-minio:9000", "pxxadmin", "pxxminio123")

	cache.InitBloomFilter(1000000, 3)
	log.Println("[Bloom] filter ready")

	// 5. 缓存管理器
	cm := cache.NewCacheManager()

	// 6. 初始化 Kafka
	mq.InitKafka(cfg.KafkaBros)
	log.Println("[Kafka] producer ready")

	// Start Kafka consumers
	mq.StartConsumer(cfg.KafkaBros, mq.TopicPostCreated, "pxx-group", func(data []byte) error {
		log.Printf("[Kafka] post.created received: %s", string(data))
		return nil
	})
	mq.StartConsumer(cfg.KafkaBros, mq.TopicOrderCreated, "pxx-group", func(data []byte) error {
		log.Printf("[Kafka] order.created received: %s", string(data))
		return nil
	})
	mq.StartConsumer(cfg.KafkaBros, mq.TopicGroupBuyFlash, "pxx-group", func(data []byte) error {
		log.Printf("[Kafka] groupbuy.flash received: %s", string(data))
		return nil
	})
	log.Println("[Kafka] consumers started")

	// 7. 初始化 JWT
	middleware.InitAuth(cfg.JWTSecret)

	// 8. 组装 Service -> Handler -> Router
	services := service.NewServices(db, cm)
	handlers := handler.NewHandlers(services)
	r := router.Setup(handlers)

	// 9. 启动
	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	// Dutch Auction engine: demo goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			var products []model.Product
			db.Where("status = 1 AND original_price > 0").Find(&products)
			for _, p := range products {
				newPrice := p.Price - int64(float64(p.Price)*0.05)
				if newPrice > 0 {
					db.Model(&p).Update("price", newPrice)
					cm.Invalidate(context.Background(), cache.BuildKey("product", p.ID))
				}
			}
			log.Printf("[Auction] tick: %d products checked", len(products))
		}
	}()


	// Credit score engine: update every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			var users []model.User
			db.Find(&users)
			for _, u := range users {
				// Base score 60 + dispute losses
				score := 60
				// Count dispute losses
				var lossCount int64
				db.Model(&model.DisputeCase{}).Where("(plaintiff_id = ? OR defendant_id = ?) AND winner != '' AND status = 2", u.ID, u.ID).Count(&lossCount)
				// Each loss: -5 points
				score -= int(lossCount) * 5
				if score < 0 { score = 0 }
				if score > 100 { score = 100 }
				db.Model(&u).Update("credit", score)
			}
			log.Printf("[Credit] engine tick: %d users updated", len(users))
		}
	}()

	log.Printf("[Server] starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[Server] failed: %v", err)
	}
}
