package handler

import (
	"context"

	"pxx/feed/service"
)

// FeedHandler 预留的 Feed gRPC 处理器
// Feed 服务当前主要通过 HTTP Gateway 访问或 Kafka 异步消费
// 未来可扩展 gRPC 接口用于 Feed 服务间通信
type FeedHandler struct {
	postSvc *service.PostService
}

func NewFeedHandler(postSvc *service.PostService) *FeedHandler {
	return &FeedHandler{postSvc: postSvc}
}
