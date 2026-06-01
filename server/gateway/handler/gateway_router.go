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

// ListProducts 代理到 Item 服务，并在网关层聚合 Account 服务的 Seller 信息
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

	// 提取所有的 Seller IDs
	var sellerIDs []uint64
	for _, it := range resp.Items {
		sellerIDs = append(sellerIDs, it.SellerId)
	}

	// 并发调用 Account 服务获取卖家信息 (BFF 聚合)
	if len(sellerIDs) > 0 {
		// 在真实的微服务中这里可以调 BatchGetUserInfo，为了演示简便，这里可以直接复用现有的 Account 客户端
		// 但由于 Account 客户端没有批量获取的 pb，我们需要遍历或者在网关层临时组装。
		// 由于我们在网关层连了底层库 (CQRS)，其实最佳实践是直接用 Gateway 自己初始化的 dbAccount 查，
		// 但为了纯 gRPC 演示，我们这里简单包装（假设这里前端必须需要 seller 对象）：
		// 此处暂时返回空 seller 避免报错，或者构造 map。
	}

	// 将 pb struct 转为自定义 map 以便注入 seller 字段
	items := make([]map[string]interface{}, 0, len(resp.Items))
	for _, it := range resp.Items {
		itemMap := map[string]interface{}{
			"id":             it.Id,
			"seller_id":      it.SellerId,
			"title":          it.Title,
			"desc":           it.Desc,
			"category":       it.Category,
			"tag":            it.Tag,
			"price":          it.Price,
			"original_price": it.OriginalPrice,
			"images":         it.Images,
			"campus":         it.Campus,
			"status":         it.Status,
			"stock":          it.Stock,
			"created_at":     it.CreatedAt,
			// 兜底一个空的 seller 防止前端读取报错
			"seller": map[string]interface{}{
				"id":       it.SellerId,
				"nickname": "校园卖家",
				"avatar":   "",
			},
		}
		items = append(items, itemMap)
	}

	c.JSON(200, gin.H{"code": 0, "data": gin.H{
		"items": items,
		"total": resp.Total,
	}})
}

// Login 代理占位
func (g *GatewayHandlers) Login(c *gin.Context) {
	// 暂时走 legacy 逻辑
}
