package handler

import (
	"strconv"

	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	svc *service.RentalService
}

func NewRentalHandler(svc *service.RentalService) *RentalHandler {
	return &RentalHandler{svc: svc}
}

func (h *RentalHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, gin.H{"total": total, "items": items})
}

func (h *RentalHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Detail(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, 404, "租赁物品不存在")
		return
	}
	response.OK(c, item)
}

func (h *RentalHandler) CreateOrder(c *gin.Context) {
	rentalID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var order model.RentalOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	order.RentalID = uint(rentalID)
	order.Status = model.RentalStatusPending
	if err := h.svc.CreateOrder(c.Request.Context(), &order); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, order)
}

func (h *RentalHandler) Return(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Return(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, nil)
}
