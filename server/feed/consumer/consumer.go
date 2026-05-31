package consumer

import (
	"context"
	"log"

	"pxx/feed/service"
	"pxx/internal/mq"
)

// StartFeedConsumer 启动 Feed 流 Kafka 消费者
func StartFeedConsumer(ctx context.Context, brokers string, svc *service.PostService) {
	mq.StartConsumer(ctx, brokers, mq.TopicFeedEvent, "pxx-feed-consumer", func(data []byte) error {
		if err := svc.ProcessEvent(data); err != nil {
			log.Printf("[FeedConsumer] process error: %v", err)
			return err
		}
		log.Printf("[FeedConsumer] processed event: %s", string(data[:min(len(data), 100)]))
		return nil
	})
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
