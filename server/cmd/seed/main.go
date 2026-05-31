package main

import (
	"log"
	"time"

	accountModel "pxx/account/model"
	feedModel "pxx/feed/model"
	itemModel "pxx/item/model"
	"pxx/pkg/config"
	"pxx/pkg/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 1. Connect Account DB
	accountDSN := cfg.AccountDB.DSN()
	accountDB, err := gorm.Open(mysql.Open(accountDSN), &gorm.Config{})
	if err != nil { log.Fatalf("[Seed] account db connect failed: %v", err) }
	accountDB.Exec("SET NAMES utf8mb4")

	// 2. Connect Item DB
	itemDSN := cfg.ItemDB.DSN()
	itemDB, err := gorm.Open(mysql.Open(itemDSN), &gorm.Config{})
	if err != nil { log.Fatalf("[Seed] item db connect failed: %v", err) }
	itemDB.Exec("SET NAMES utf8mb4")

	// 3. Connect Feed DB
	feedDSN := cfg.FeedDB.DSN()
	feedDB, err := gorm.Open(mysql.Open(feedDSN), &gorm.Config{})
	if err != nil { log.Fatalf("[Seed] feed db connect failed: %v", err) }
	feedDB.Exec("SET NAMES utf8mb4")

	log.Println("[Seed] databases connected, starting seed...")

	seedAccount(accountDB)
	seedItem(itemDB)
	seedExtendedItem(itemDB)
	seedFeed(feedDB)

	log.Println("[Seed] done.")
}

func seedAccount(db *gorm.DB) {
	db.AutoMigrate(&accountModel.User{})

	u1 := accountModel.User{OpenID: "test_001", Nickname: "校园小明", Campus: "南校区", Department: "计算机学院", Balance: 100000}
	u2 := accountModel.User{OpenID: "test_002", Nickname: "图书馆小张", Campus: "北校区", Department: "软件学院", Balance: 50000}
	u3 := accountModel.User{OpenID: "admin_001", Nickname: "校园评审员", Campus: "管理", Department: "信息中心", Role: "admin"}
	
	u1.ID = 1; u2.ID = 2; u3.ID = 3

	db.Where(&accountModel.User{OpenID: u1.OpenID}).FirstOrCreate(&u1)
	db.Where(&accountModel.User{OpenID: u2.OpenID}).FirstOrCreate(&u2)
	db.Where(&accountModel.User{OpenID: u3.OpenID}).FirstOrCreate(&u3)
	log.Println("[Seed Account] inserted or verified users")
}

func seedItem(db *gorm.DB) {
	db.AutoMigrate(&itemModel.Product{})

	p1 := itemModel.Product{SellerID: 1, Title: "考研408全套资料", Desc: "数据结构+组成原理+操作系统+计算机网络，几乎全新", Category: "考研资料", Tag: "计算机", Price: 4500, OriginalPrice: 12000, Campus: "南校区", Stock: 1, Status: 1}
	p2 := itemModel.Product{SellerID: 2, Title: "罗技G502鼠标", Desc: "用了半年，功能完好，箱说全", Category: "数码外设", Tag: "鼠标", Price: 15000, OriginalPrice: 39900, Campus: "北校区", Stock: 2, Status: 1}
	p3 := itemModel.Product{SellerID: 1, Title: "宿舍用小台灯", Desc: "LED台灯，三档调光", Category: "日用百货", Tag: "台灯", Price: 2500, OriginalPrice: 5000, Campus: "南校区", Stock: 5, Status: 1}

	db.Where(&itemModel.Product{Title: p1.Title}).FirstOrCreate(&p1)
	db.Where(&itemModel.Product{Title: p2.Title}).FirstOrCreate(&p2)
	db.Where(&itemModel.Product{Title: p3.Title}).FirstOrCreate(&p3)
	log.Println("[Seed Item] inserted or verified products")
}

func seedExtendedItem(db *gorm.DB) {
	db.AutoMigrate(&itemModel.GroupBuy{}, &itemModel.RentalItem{}, &itemModel.BarterItem{}, &itemModel.Subscription{})

	now := time.Now()
	g1 := itemModel.GroupBuy{Title: "瑞幸生椰拿铁10人团", Desc: "食堂门口瑞幸自提", PricePerUnit: 990, MinPeople: 10, CurrentPeople: 1, TotalStock: 50, SoldCount: 1, StartAt: now, EndAt: now.Add(24 * time.Hour), Status: 2}
	r1 := itemModel.RentalItem{OwnerID: 1, Title: "大疆单反相机出租", Desc: "适合周末出去玩，带两块电池", DailyPrice: 5000, Deposit: 200000, Status: 1}
	b1 := itemModel.BarterItem{UserID: 2, Title: "全新机械键盘换降噪耳机", Desc: "刚买的红轴机械键盘，想换一个成色还可以的降噪耳机", WantItem: "降噪耳机", Status: 1}
	
	actEnc, _ := utils.AESEncrypt("netflix_test_account@example.com")
	pwdEnc, _ := utils.AESEncrypt("netflix_secret_pass")
	s1 := itemModel.Subscription{OwnerID: 1, Brand: "Netflix", Title: "Netflix 4K高级会员拼车", Desc: "长期稳定车队", PricePerDay: 200, PricePerWeek: 1200, AccountCipher: actEnc, PassCipher: pwdEnc, Status: 1}

	db.Where(&itemModel.GroupBuy{Title: g1.Title}).FirstOrCreate(&g1)
	db.Where(&itemModel.RentalItem{Title: r1.Title}).FirstOrCreate(&r1)
	db.Where(&itemModel.BarterItem{Title: b1.Title}).FirstOrCreate(&b1)
	db.Where(&itemModel.Subscription{Title: s1.Title}).FirstOrCreate(&s1)
	
	log.Println("[Seed Item] inserted extended items (groupbuy, rental, barter, sub)")
}

func seedFeed(db *gorm.DB) {
	db.AutoMigrate(&feedModel.Post{}, &feedModel.Task{}, &feedModel.DisputeCase{})

	post1 := feedModel.Post{UserID: 1, Title: "食堂探店", Content: "学校南门新开的烤肉店太好吃了，强烈推荐！", Type: 2, Tags: "美食,日常", Status: 1}
	task1 := feedModel.Task{PublisherID: 1, Title: "急求代拿快递", Desc: "南门韵达快递，大件，送到男生宿舍5栋", FromPlace: "南门韵达", ToPlace: "男生宿舍5栋", Reward: 800, Status: 1}
	dispute1 := feedModel.DisputeCase{Title: "买到的二手教材有大面积涂鸦", PlaintiffID: 1, DefendantID: 2, OrderID: 1, PlaintiffDesc: "卖家说全新，结果全是笔记", DefendantDesc: "笔记不影响阅读啊", Status: 1}

	db.Where(&feedModel.Post{Content: post1.Content}).FirstOrCreate(&post1)
	db.Where(&feedModel.Task{Title: task1.Title}).FirstOrCreate(&task1)
	db.Where(&feedModel.DisputeCase{Title: dispute1.Title}).FirstOrCreate(&dispute1)
	
	log.Println("[Seed Feed] inserted posts, tasks, and disputes")
}