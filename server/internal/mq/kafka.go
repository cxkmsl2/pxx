package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitKafka(brokers string) {
	writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	}
	log.Println("[Kafka] producer ready")
}

func GetWriter() *kafka.Writer {
	return writer
}

// Publish 发送消息到指定 Topic
func Publish(ctx context.Context, topic string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: data,
	})
}

// 本项目的 Kafka Topics
const (
	TopicPostCreated   = "pxx.post.created"   // 帖子发布 → 图片压缩/审核
	TopicOrderCreated  = "pxx.order.created"   // 订单创建 → 异步处理
	TopicOrderExpired  = "pxx.order.expired"   // 延迟消费：15 分钟未支付取消
	TopicGroupBuyFlash = "pxx.groupbuy.flash"  // 拼团秒杀削峰
	TopicFeedEvent     = "pxx.feed.event"       // Feed 流事件
)

// StartConsumer 启动一个带断线重连和 Context 控制的消费者
func StartConsumer(ctx context.Context, brokers, topic, groupID string, handler func([]byte) error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := func() error {
					r := kafka.NewReader(kafka.ReaderConfig{
						Brokers:  []string{brokers},
						Topic:    topic,
						GroupID:  groupID,
						MinBytes: 10e3,
						MaxBytes: 10e6,
					})
					defer r.Close()

					log.Printf("[Kafka] consumer started: topic=%s group=%s", topic, groupID)

					for {
						// 使用 ctx 读取消息，以便在外部 cancel 时能立即退出
						msg, err := r.ReadMessage(ctx)
						if err != nil {
							return err
						}
						if err := handler(msg.Value); err != nil {
							log.Printf("[Kafka] handler error: %v", err)
						}
					}
				}()

				if err != nil {
					select {
					case <-ctx.Done():
						return
					default:
						log.Printf("[Kafka] consumer error (topic=%s group=%s): %v, reconnecting in 3s...", topic, groupID, err)
						time.Sleep(3 * time.Second)
					}
				}
			}
		}
	}()
}
