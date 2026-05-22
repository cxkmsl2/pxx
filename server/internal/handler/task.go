package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) List(c *gin.Context) {
	taskType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.List(c.Request.Context(), taskType, page, pageSize)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, gin.H{"total": total, "items": items})
}

func (h *TaskHandler) Create(c *gin.Context) {
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	t.PublisherID = middleware.GetUserID(c)
	if err := h.svc.Create(c.Request.Context(), &t); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, t)
}

func (h *TaskHandler) Take(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	takerID := middleware.GetUserID(c)
	if err := h.svc.Take(c.Request.Context(), uint(id), takerID); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, nil)
}

func (h *TaskHandler) Done(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	if err := h.svc.Done(c.Request.Context(), uint(id), userID); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, nil)
}
