package handler

import (
	"pxx/gateway/client"
	"pxx/internal/handler"
	"pxx/internal/service"

	"gorm.io/gorm"
)

// GatewayHandlers 网关层 Handler
type GatewayHandlers struct {
	clients *client.Clients
}

func NewGatewayHandlers(clients *client.Clients) *GatewayHandlers {
	return &GatewayHandlers{clients: clients}
}

// InitHandlers 初始化真正的 Handler（使用真实数据库连接）
func InitHandlers(db *gorm.DBByService, cm interface{}) *handler.Handlers {
	// 需要引用 internal/cache 的类型
	return handler.NewHandlers(nil) // placeholder
}

func (g *GatewayHandlers) SetDb(db *gorm.DB) *handler.Handlers {
	// 过渡方案：Gateway 连接旧的 pxx 数据库
	// 使用 internal/service 的完整实现
	_ = db
	return nil
}
