package service

import (
	"pxx/internal/cache"

	"gorm.io/gorm"
)

type Services struct {
	User         *UserService
	Product      *ProductService
	Order        *OrderService
	Post         *PostService
	GroupBuy     *GroupBuyService
	Rental       *RentalService
	Barter       *BarterService
	Subscription *SubscriptionService
	Task         *TaskService
	Message      *MessageService
	Dispute      *DisputeService
}

func NewServices(accountDB, itemDB, tradeDB, feedDB *gorm.DB, cm *cache.CacheManager) *Services {
	svc := &Services{}
	
	// Initialize base services first without dependencies
	svc.User = NewUserService(accountDB)
	svc.Product = NewProductService(itemDB, cm)
	svc.Order = NewOrderService(tradeDB)
	svc.Post = NewPostService(feedDB, cm)
	svc.GroupBuy = NewGroupBuyService(itemDB, cm)
	svc.Rental = NewRentalService(itemDB, cm)
	svc.Barter = NewBarterService(itemDB, cm)
	svc.Subscription = NewSubscriptionService(itemDB, cm)
	svc.Task = NewTaskService(feedDB, cm)
	svc.Message = NewMessageService(feedDB)
	svc.Dispute = NewDisputeService(feedDB)
	
	// Now inject cross-service dependencies
	svc.Message.UserService = svc.User
	svc.Order.ProductService = svc.Product
	svc.Product.UserService = svc.User
	svc.Rental.UserService = svc.User
	svc.Subscription.UserService = svc.User
	
	return svc
}
