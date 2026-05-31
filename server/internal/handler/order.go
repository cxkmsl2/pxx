package handler

import (
	"fmt"
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc    *service.OrderService
	Create func(c *gin.Context)
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	h := &OrderHandler{svc: svc}
	h.Create = h.create
	return h
}

func (h *OrderHandler) create(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	buyerID := middleware.GetUserID(c)
	order, err := h.svc.Create(c.Request.Context(), buyerID, req.ProductID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, order)
}

func (h *OrderHandler) ListMy(c *gin.Context) {
	userID := middleware.GetUserID(c)
	fmt.Printf("[DEBUG] ListMy called for userID: %d\n", userID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orders, total, err := h.svc.ListMy(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	fmt.Printf("[DEBUG] ListMy found %d orders\n", total)
	response.OK(c, gin.H{"total": total, "items": orders})
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, nil)
}
