package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) List(c *gin.Context) {
	ptype, _ := strconv.Atoi(c.Query("type"))
	tag := c.Query("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	posts, total, err := h.svc.List(c.Request.Context(), int8(ptype), tag, page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, gin.H{"total": total, "items": posts})
}

func (h *PostHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Detail(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, 404, "帖子不存在")
		return
	}
	response.OK(c, p)
}

func (h *PostHandler) Create(c *gin.Context) {
	var p model.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	p.UserID = middleware.GetUserID(c)
	if err := h.svc.Create(c.Request.Context(), &p); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, p)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.svc.Update(c.Request.Context(), uint(id), updates); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, nil)
}
