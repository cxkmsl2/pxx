package service

import (
	"context"
	"time"

	"pxx/internal/cache"
	"pxx/internal/model"
	"pxx/pkg/utils"

	"gorm.io/gorm"
)

type ProductService struct {
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewProductService(db *gorm.DB, cm *cache.CacheManager) *ProductService {
	return &ProductService{db: db, cache: cm}
}

func (s *ProductService) List(ctx context.Context, category, campus, keyword string, sellerID uint, page, pageSize int) ([]model.Product, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var products []model.Product

	query := s.db.WithContext(ctx).Model(&model.Product{}).Where("status = 1")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if campus != "" {
		query = query.Where("campus LIKE ?", "%"+campus+"%")
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR `desc` LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sellerID > 0 {
		query = query.Where("seller_id = ?", sellerID)
	}

	query.Count(&total)
	query.Preload("Seller").Order("created_at DESC").Offset(offset).Limit(limit).Find(&products)
	return products, total, nil
}

func (s *ProductService) Detail(ctx context.Context, id uint) (*model.Product, error) {
	key := cache.BuildKey("product", id)
	var p model.Product
	err := s.cache.GetOrLoad(ctx, key, &p, 30*time.Minute, func(ctx context.Context) (interface{}, error) {
		var product model.Product
		if err := s.db.WithContext(ctx).Preload("Seller").First(&product, id).Error; err != nil {
			return nil, err
		}
		// 更新浏览量
		s.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).
			UpdateColumn("view_count", gorm.Expr("view_count + 1"))
		return &product, nil
	})
	return &p, err
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	err := s.db.WithContext(ctx).Create(p).Error
	if err == nil {
		cache.GetBloomFilter().Add(cache.BuildKey("product", p.ID))
	}
	return err
}

func (s *ProductService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	err := s.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Updates(updates).Error
	if err == nil {
		s.cache.Invalidate(ctx, cache.BuildKey("product", id))
	}
	return err
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Update("status", 0).Error
	if err == nil {
		s.cache.Invalidate(ctx, cache.BuildKey("product", id))
	}
	return err
}
