package handler

import (
	"pxx/internal/middleware"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		OpenID string `json:"open_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	result, err := h.svc.Login(c.Request.Context(), req.OpenID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, gin.H{
		"user":  result.User,
		"token": result.Token,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.svc.Update(c.Request.Context(), userID, updates); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, nil)
}
