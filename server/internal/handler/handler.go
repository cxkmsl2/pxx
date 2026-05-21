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
	Barter   *BarterHandler
}

func NewHandlers(svc *service.Services) *Handlers {
	return &Handlers{
		User:     NewUserHandler(svc.User),
		Product:  NewProductHandler(svc.Product),
		Order:    NewOrderHandler(svc.Order),
		Post:     NewPostHandler(svc.Post),
		GroupBuy: NewGroupBuyHandler(svc.GroupBuy),
		Rental:   NewRentalHandler(svc.Rental),
		Barter:   NewBarterHandler(svc.Barter),
	}
}
