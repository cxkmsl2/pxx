package service

import (
	"pxx/internal/cache"

	"gorm.io/gorm"
)

type Services struct {
	User     *UserService
	Product  *ProductService
	Order    *OrderService
	Post     *PostService
	GroupBuy *GroupBuyService
	Rental   *RentalService
	Barter       *BarterService
	Subscription *SubscriptionService
	Task         *TaskService
	Message      *MessageService
	Dispute      *DisputeService
}

func NewServices(db *gorm.DB, cm *cache.CacheManager) *Services {
	return &Services{
		User:     NewUserService(db),
		Product:  NewProductService(db, cm),
		Order:    NewOrderService(db),
		Post:     NewPostService(db, cm),
		GroupBuy: NewGroupBuyService(db, cm),
		Rental:   NewRentalService(db, cm),
		Barter:       NewBarterService(db, cm),
		Subscription: NewSubscriptionService(db, cm),
		Task:         NewTaskService(db, cm),
		Message:      NewMessageService(db),
		Dispute:      NewDisputeService(db),
	}
}
