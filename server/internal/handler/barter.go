package handler

import (
	"strconv"
	"time"

	"pxx/internal/middleware"

	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type BarterHandler struct {
	svc *service.BarterService
}

func NewBarterHandler(svc *service.BarterService) *BarterHandler {
	return &BarterHandler{svc: svc}
}

func (h *BarterHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, gin.H{"total": total, "items": items})
}

func (h *BarterHandler) Create(c *gin.Context) {
	var item model.BarterItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.svc.Create(c.Request.Context(), &item); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, item)
}

func (h *BarterHandler) Agree(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		MeetPlace string `json:"meet_place"`
		MeetTime  string `json:"meet_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	meetTime := time.Now().Add(24 * time.Hour)
	if req.MeetTime != "" {
		// Try MySQL format first, then ISO
		if t, err := time.Parse("2006-01-02 15:04:05", req.MeetTime); err == nil {
			meetTime = t
		} else if t, err := time.Parse(time.RFC3339, req.MeetTime); err == nil {
			meetTime = t
		}
	}
	userID := middleware.GetUserID(c)
	order, err := h.svc.Agree(c.Request.Context(), userID, uint(id), req.MeetPlace, meetTime)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, order)
}
