package router

import (
	"pxx/internal/handler"
	"pxx/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(h *handler.Handlers) *gin.Engine {
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())
	g.Use(middleware.CORS())

	api := g.Group("/api/v1")
	{
		// 公开接口
		api.POST("/login", h.User.Login)
		api.GET("/products", h.Product.List)
		api.GET("/products/:id", h.Product.Detail)
		api.GET("/posts", h.Post.List)
		api.GET("/posts/:id", h.Post.Detail)
		api.GET("/groupbuys", h.GroupBuy.List)
		api.GET("/groupbuys/:id", h.GroupBuy.Detail)
		api.GET("/rentals", h.Rental.List)
		api.GET("/rentals/:id", h.Rental.Detail)
		api.GET("/barters", h.Barter.List)

		// 需要认证
		auth := api.Group("", middleware.AuthRequired())
		{
			// 用户
			auth.GET("/user/profile", h.User.Profile)
			auth.PUT("/user/profile", h.User.UpdateProfile)

			// 商品
			auth.POST("/products", h.Product.Create)
			auth.PUT("/products/:id", h.Product.Update)
			auth.DELETE("/products/:id", h.Product.Delete)

			// 订单
			auth.POST("/orders", h.Order.Create)
			auth.GET("/orders", h.Order.ListMy)
			auth.PUT("/orders/:id/status", h.Order.UpdateStatus)

			// 帖子
			auth.POST("/posts", h.Post.Create)
			auth.PUT("/posts/:id", h.Post.Update)

			// 拼团
			auth.POST("/groupbuys/:id/join", h.GroupBuy.Join)

			// 租赁
			auth.POST("/rentals/:id/order", h.Rental.CreateOrder)
			auth.PUT("/rental-orders/:id/return", h.Rental.Return)

			// 以物换物
			auth.POST("/barters", h.Barter.Create)
			auth.POST("/barters/:id/agree", h.Barter.Agree)
		}
	}

	return g
}
