package router

import (
	"pxx/internal/handler"
	"pxx/internal/middleware"
	"pxx/internal/ws"

	"github.com/gin-gonic/gin"
)

func Setup(h *handler.Handlers) *gin.Engine {
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())
	g.Use(middleware.CORS())

	// 健康检查端点（供 Docker healthcheck）
	g.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := g.Group("/api/v1")
	{
		api.POST("/login", h.User.Login)
		api.GET("/users/:id", h.User.PublicProfile)
		if h.Product.List != nil {
			api.GET("/products", h.Product.List)
		}
		api.GET("/products/:id", h.Product.Detail)
		api.GET("/posts", h.Post.List)
		api.GET("/posts/:id", h.Post.Detail)
		api.GET("/groupbuys", h.GroupBuy.List)
		api.GET("/groupbuys/:id", h.GroupBuy.Detail)
		api.GET("/rentals", h.Rental.List)
		api.GET("/rentals/:id", h.Rental.Detail)
		api.GET("/barters", h.Barter.List)
		api.GET("/subscriptions", h.Subscription.List)
		api.GET("/tasks", h.Task.List)
		api.GET("/disputes", h.Dispute.List)

		auth := api.Group("", middleware.AuthRequired())
		{
			auth.GET("/user/profile", h.User.Profile)
			auth.PUT("/user/profile", h.User.UpdateProfile)

			auth.POST("/products", h.Product.Create)
			auth.PUT("/products/:id", h.Product.Update)
			auth.DELETE("/products/:id", h.Product.Delete)

			// 只有在 Create 为空时才注册（或者直接信任传入的 h.Order.Create）
			if h.Order.Create != nil {
				auth.POST("/orders", h.Order.Create)
			}
			auth.GET("/orders", h.Order.ListMy)
			auth.PUT("/orders/:id/status", h.Order.UpdateStatus)

			auth.POST("/posts", h.Post.Create)
			auth.PUT("/posts/:id", h.Post.Update)

			auth.POST("/groupbuys/:id/join", h.GroupBuy.Join)

			auth.POST("/rentals/:id/order", h.Rental.CreateOrder)
			auth.PUT("/rental-orders/:id/return", h.Rental.Return)

			auth.POST("/barters", h.Barter.Create)
			auth.POST("/barters/:id/agree", h.Barter.Agree)

			auth.POST("/subscriptions", h.Subscription.Create)
			auth.POST("/subscriptions/:id/rent", h.Subscription.Rent)
			auth.GET("/subscriptions/orders/:id/credential", h.Subscription.Credential)
			auth.POST("/subscriptions/orders/:id/dispute", h.Subscription.Dispute)

			auth.POST("/tasks", h.Task.Create)
			auth.POST("/tasks/:id/take", h.Task.Take)
			auth.POST("/tasks/:id/done", h.Task.Done)

			auth.POST("/disputes", h.Dispute.Create)
			auth.POST("/disputes/:id/vote", h.Dispute.Vote)
			auth.POST("/disputes/:id/accept", h.Dispute.Accept)
			auth.POST("/upload", handler.UploadImage)

			auth.POST("/messages", h.Message.Send)
			auth.GET("/messages/sys", h.Message.SysMessages)
			auth.GET("/messages/conversations", h.Message.Conversations)
			auth.GET("/messages/:user_id", h.Message.GetMessages)
			auth.POST("/messages/read/:user_id", h.Message.MarkRead)
		}
	}

	// WebSocket
	wsGroup := g.Group("/api/v1/ws")
	wsGroup.Use(middleware.AuthRequired())
	wsGroup.GET("/connect", func(c *gin.Context) {
		ws.HandleWS(c)
	})

	g.Static("/api/v1/uploads", "./static/uploads")
	return g
}