package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
}

func NewSubscriptionHandler(svc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	brand := c.Query("brand")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.List(c.Request.Context(), brand, page, pageSize)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, gin.H{"total": total, "items": items})
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	// AccountCipher and PassCipher come as plaintext from frontend, will be encrypted in service
	// The frontend sends "account" and "password" fields, map them
	var req struct {
		Brand       string `json:"brand"`
		Title       string `json:"title"`
		Desc        string `json:"desc"`
		Account     string `json:"account"`
		Password    string `json:"password"`
		PricePerDay int64  `json:"price_per_day"`
		PricePerWeek int64 `json:"price_per_week"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	sub = model.Subscription{
		OwnerID:      middleware.GetUserID(c),
		Brand:        req.Brand,
		Title:        req.Title,
		Desc:         req.Desc,
		AccountCipher: req.Account,
		PassCipher:   req.Password,
		PricePerDay:  req.PricePerDay,
		PricePerWeek: req.PricePerWeek,
	}
	if err := h.svc.Create(c.Request.Context(), &sub); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, sub)
}

func (h *SubscriptionHandler) Rent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Days int `json:"days" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	buyerID := middleware.GetUserID(c)
	order, err := h.svc.Rent(c.Request.Context(), uint(id), buyerID, req.Days)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, order)
}

func (h *SubscriptionHandler) Credential(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	account, password, rentEnd, err := h.svc.GetCredential(c.Request.Context(), uint(id), userID)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, gin.H{"account": account, "password": password, "rent_end": rentEnd})
}

func (h *SubscriptionHandler) Dispute(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	userID := middleware.GetUserID(c)
	if err := h.svc.Dispute(c.Request.Context(), uint(id), userID, req.Reason); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, nil)
}
