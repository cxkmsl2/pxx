package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type DisputeHandler struct {
	svc *service.DisputeService
}

func NewDisputeHandler(svc *service.DisputeService) *DisputeHandler {
	return &DisputeHandler{svc: svc}
}

func (h *DisputeHandler) List(c *gin.Context) {
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cases, total, err := h.svc.ListCases(c.Request.Context(), int8(status), page, pageSize)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, gin.H{"total": total, "items": cases})
}

func (h *DisputeHandler) Create(c *gin.Context) {
	var cReq model.DisputeCase
	if err := c.ShouldBindJSON(&cReq); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	cReq.PlaintiffID = middleware.GetUserID(c)
	if err := h.svc.Create(c.Request.Context(), &cReq); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, cReq)
}

func (h *DisputeHandler) Accept(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Accept(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err.Error()); return
	}
	response.OK(c, nil)
}

func (h *DisputeHandler) Vote(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		VoteFor string `json:"vote_for" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"); return
	}
	voterID := middleware.GetUserID(c)
	result, err := h.svc.Vote(c.Request.Context(), uint(id), voterID, req.VoteFor)
	if err != nil { response.Error(c, err.Error()); return }
	response.OK(c, gin.H{"result": result})
}
