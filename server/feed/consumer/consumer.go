package consumer

import (
	"context"
	"log"

	"pxx/feed/service"

	"github.com/segmentio/kafka-go"
)

// StartFeedConsumer 启动 Feed 流 Kafka 消费者
// 消费 pxx.feed.event 主题，批量写库
func StartFeedConsumer(brokers string, svc *service.PostService) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokers},
		Topic:     "pxx.feed.event",
		GroupID:   "pxx-feed-consumer",
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	go func() {
		defer reader.Close()
		log.Println("[FeedConsumer] started, consuming pxx.feed.event")

		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[FeedConsumer] read error: %v", err)
				continue
			}

			if err := svc.ProcessEvent(msg.Value); err != nil {
				log.Printf("[FeedConsumer] process error: %v", err)
			} else {
				log.Printf("[FeedConsumer] processed event: %s", string(msg.Value[:min(len(msg.Value), 100)]))
			}
		}
	}()
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
