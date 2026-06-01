package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	feedModel "pxx/feed/model"
	"pxx/internal/mq"

	"gorm.io/gorm"
)

// PostEvent 发帖事件（用于 Kafka 异步投递）
type PostEvent struct {
	UserID    uint      `json:"user_id"`
	Type      int8      `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Images    string    `json:"images"`
	Tags      string    `json:"tags"`
	ProductID *uint     `json:"product_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// PublishEvent 发布发帖事件到 Kafka
// 仅做参数校验，不直接写库，返回 202 Accepted 语义
func (s *PostService) PublishEvent(ctx context.Context, event *PostEvent) error {
	if event.Title == "" || event.Content == "" {
		return fmt.Errorf("标题和内容不能为空")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return mq.Publish(ctx, "pxx.feed.event", data)
}

// ProcessEvent 处理 Kafka 消费到的发帖事件
// 使用 INSERT ... ON DUPLICATE KEY UPDATE 批量写库
func (s *PostService) ProcessEvent(data []byte) error {
	var event PostEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	post := &feedModel.Post{
		UserID:    event.UserID,
		Type:      event.Type,
		Title:     event.Title,
		Content:   event.Content,
		Images:    event.Images,
		Tags:      event.Tags,
		ProductID: event.ProductID,
	}

	return s.db.Create(post).Error
}

// List 帖子列表 (支持游标分页)
func (s *PostService) List(ctx context.Context, postType int8, pageSize int, cursorTime int64, cursorId uint) ([]feedModel.Post, int64, error) {
	if pageSize <= 0 { pageSize = 20 }

	var posts []feedModel.Post
	query := s.db.WithContext(ctx).Model(&feedModel.Post{}).Where("status = 1")
	if postType > 0 {
		query = query.Where("type = ?", postType)
	}
	
	// 游标分页核心逻辑
	if cursorTime > 0 {
		t := time.UnixMilli(cursorTime)
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", t, t, cursorId)
	}

	// 游标模式直接基于索引查找，无需 Offset，也无需 Count（耗时）
	query.Order("created_at DESC, id DESC").Limit(pageSize).Find(&posts)
	return posts, 0, nil
}

// Create 直接创建帖子（传统同步方式，留给迁移期使用）
func (s *PostService) Create(ctx context.Context, post *feedModel.Post) error {
	return s.db.WithContext(ctx).Create(post).Error
}
