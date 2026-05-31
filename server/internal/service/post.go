package service

import (
	"context"
	"time"

	"pxx/internal/cache"
	"pxx/internal/model"
	"pxx/internal/mq"
	"pxx/pkg/utils"

	"gorm.io/gorm"
)

type PostService struct {
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewPostService(db *gorm.DB, cm *cache.CacheManager) *PostService {
	return &PostService{db: db, cache: cm}
}

func (s *PostService) List(ctx context.Context, ptype int8, tag string, page, pageSize int) ([]model.Post, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var posts []model.Post
	query := s.db.WithContext(ctx).Model(&model.Post{}).Where("status = 1")
	if ptype > 0 {
		query = query.Where("type = ?", ptype)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)
	return posts, total, nil
}

func (s *PostService) Detail(ctx context.Context, id uint) (*model.Post, error) {
	// 先查存在性，再更新浏览量
	var exists int64
	s.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Count(&exists)
	if exists > 0 {
		s.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).
			UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	}

	key := cache.BuildKey("post", id)
	var p model.Post
	err := s.cache.GetOrLoad(ctx, key, &p, 15*time.Minute, func(ctx context.Context) (interface{}, error) {
		var post model.Post
		if err := s.db.WithContext(ctx).First(&post, id).Error; err != nil {
			return nil, err
		}
		return &post, nil
	})
	return &p, err
}

func (s *PostService) Create(ctx context.Context, p *model.Post) error {
	err := s.db.WithContext(ctx).Create(p).Error
	if err == nil {
		// 异步：发到 Kafka 处理图片压缩/审核
		_ = mq.Publish(ctx, mq.TopicPostCreated, p)
	}
	return err
}

func (s *PostService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}