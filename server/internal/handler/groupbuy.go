package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type GroupBuyHandler struct {
	svc *service.GroupBuyService
}

func NewGroupBuyHandler(svc *service.GroupBuyService) *GroupBuyHandler {
	return &GroupBuyHandler{svc: svc}
}

func (h *GroupBuyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, gin.H{"total": total, "items": items})
}

func (h *GroupBuyHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Detail(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, 404, "拼团不存在")
		return
	}
	response.OK(c, item)
}

func (h *GroupBuyHandler) Join(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Quantity int  `json:"quantity"`
		IsLeader bool `json:"is_leader"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if req.Quantity == 0 {
		req.Quantity = 1
	}
	userID := middleware.GetUserID(c)
	if err := h.svc.Join(c.Request.Context(), uint(id), userID, req.Quantity, req.IsLeader); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, nil)
}
