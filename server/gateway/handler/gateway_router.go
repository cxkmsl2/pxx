package handler

import (
	"pxx/gateway/client"
	"pxx/idl/gen/item"
	"pxx/idl/gen/trade"

	"github.com/gin-gonic/gin"
)

type GatewayHandlers struct {
	clients *client.Clients
}

func NewGatewayHandlers(clients *client.Clients) *GatewayHandlers {
	return &GatewayHandlers{clients: clients}
}

// CreateOrder 代理到 Trade 服务 (TCC 事务入口)
func (g *GatewayHandlers) CreateOrder(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := uint64(val.(uint))

	var req struct {
		ProductID uint64 `json:"product_id"`
		Amount    int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// 为了防止前端不传价格，网关层主动去 Item 服务查询价格
	if req.Amount == 0 {
		productResp, err := g.clients.Item.GetProductDetail(c.Request.Context(), &item.GetProductDetailReq{ProductId: req.ProductID})
		if err == nil && productResp != nil {
			req.Amount = productResp.Price
		}
	}

	resp, err := g.clients.Trade.CreateOrder(c.Request.Context(), &trade.CreateOrderReq{
		BuyerId:   userID,
		ProductId: req.ProductID,
		Amount:    req.Amount,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	// 这里要注意，如果 TCC 阶段 Try 失败，微服务虽然不会抛 err，但是会返回状态非 SUCCESS
	if resp.Status != 1 {
		c.JSON(200, gin.H{"code": 400, "msg": resp.Msg, "data": resp})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": resp})
}

// ListProducts 代理到 Item 服务
func (g *GatewayHandlers) ListProducts(c *gin.Context) {
	var req item.ListProductsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		// ignore err
	}
	if req.Page <= 0 { req.Page = 1 }
	if req.PageSize <= 0 { req.PageSize = 20 }

	resp, err := g.clients.Item.ListProducts(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": gin.H{
		"items": resp.Items,
		"total": resp.Total,
	}})
}

// Login 代理占位
func (g *GatewayHandlers) Login(c *gin.Context) {
	// 暂时走 legacy 逻辑
}
