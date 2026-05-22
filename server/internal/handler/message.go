package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	svc *service.MessageService
}

func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) Send(c *gin.Context) {
	var req struct {
		ToUserID uint   `json:"to_user_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
		MsgType  string `json:"msg_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	if req.MsgType == "" { req.MsgType = "text" }
	fromID := middleware.GetUserID(c)
	msg, err := h.svc.Send(c.Request.Context(), fromID, req.ToUserID, req.Content, req.MsgType)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, msg)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	targetID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	userID := middleware.GetUserID(c)
	msgs, err := h.svc.GetMessages(c.Request.Context(), userID, uint(targetID), 50)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, msgs)
}

func (h *MessageHandler) SysMessages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	msgs, _ := h.svc.GetSystemMessages(c.Request.Context(), userID)
	response.OK(c, msgs)
}

func (h *MessageHandler) MarkRead(c *gin.Context) {
	fromID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	userID := middleware.GetUserID(c)
	h.svc.MarkRead(c.Request.Context(), userID, uint(fromID))
	response.OK(c, nil)
}

func (h *MessageHandler) Conversations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	convs, err := h.svc.GetConversations(c.Request.Context(), userID)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, convs)
}
