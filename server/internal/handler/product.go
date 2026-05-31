package handler

import (
	"strconv"

	"pxx/internal/middleware"
	"pxx/internal/model"
	"pxx/internal/service"
	"pxx/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc  *service.ProductService
	List func(c *gin.Context)
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	h := &ProductHandler{svc: svc}
	h.List = h.list
	return h
}

func (h *ProductHandler) list(c *gin.Context) {
	category := c.Query("category")
	campus := c.Query("campus")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	sellerStr := c.Query("seller_id")
	sellerID, _ := strconv.ParseUint(sellerStr, 10, 64)
	products, total, err := h.svc.List(c.Request.Context(), category, campus, keyword, uint(sellerID), page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     products,
	})
}

func (h *ProductHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Detail(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, 404, "商品不存在")
		return
	}
	response.OK(c, p)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var p model.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	p.SellerID = middleware.GetUserID(c)
	if err := h.svc.Create(c.Request.Context(), &p); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, p)
}

func (h *ProductHandler) Update(c *gin.Context) {
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

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.OK(c, nil)
}
