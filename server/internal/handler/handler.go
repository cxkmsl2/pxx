package handler

import (
	"pxx/internal/service"
)

type Handlers struct {
	User     *UserHandler
	Product  *ProductHandler
	Order    *OrderHandler
	Post     *PostHandler
	GroupBuy *GroupBuyHandler
	Rental   *RentalHandler
	Barter       *BarterHandler
	Subscription *SubscriptionHandler
	Task         *TaskHandler
	Message      *MessageHandler
	Dispute      *DisputeHandler
}

func NewHandlers(svc *service.Services) *Handlers {
	h := &Handlers{
		User:     NewUserHandler(svc.User),
		Product:  NewProductHandler(svc.Product),
		Order:    NewOrderHandler(svc.Order),
		Post:     NewPostHandler(svc.Post),
		GroupBuy: NewGroupBuyHandler(svc.GroupBuy),
		Rental:   NewRentalHandler(svc.Rental),
		Barter:   NewBarterHandler(svc.Barter),
		Subscription: NewSubscriptionHandler(svc.Subscription),
		Task:         NewTaskHandler(svc.Task),
		Message:      NewMessageHandler(svc.Message),
		Dispute:      NewDisputeHandler(svc.Dispute),
	}
	// 关键修复：显式注入 MessageHandler 所需的 UserService
	h.Message.UserService = svc.User
	return h
}
