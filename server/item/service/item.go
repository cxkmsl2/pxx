package service

import (
	"context"
	"fmt"

	itemModel "pxx/item/model"
	"pxx/idl/gen/item"
	"pxx/internal/cache"

	"gorm.io/gorm"
)

type ItemService struct {
	db *gorm.DB
}

func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{db: db}
}

func branchID(txID string, service string) uint64 {
	var hash uint64 = 5381
	for _, c := range txID + service {
		hash = ((hash << 5) + hash) + uint64(c)
	}
	return hash
}

// TryLockStock 锁定库存（TCC Try 阶段）
func (s *ItemService) TryLockStock(ctx context.Context, req *item.TryLockStockReq) (*item.TryLockStockResp, error) {
	branchID := branchID(req.TxId, "item")

	// 1. 幂等性检查
	var existing itemModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ?", req.TxId, branchID).First(&existing).Error
	if err == nil {
		return &item.TryLockStockResp{
			Success:  existing.Status == "LOCKED" || existing.Status == "CONFIRMED",
			BranchId: branchID,
			Msg:      "幂等返回",
		}, nil
	}

	// 2. 先写日志，防悬挂
	log := &itemModel.TCCLog{
		TxID:      req.TxId,
		BranchID:  branchID,
		ActionType: "LOCK_STOCK",
		Status:    "TRYING",
		ProductID: req.ProductId,
	}
	if err := s.db.WithContext(ctx).Create(log).Error; err != nil {
		return &item.TryLockStockResp{Success: false, Msg: "日志写入失败"}, nil
	}

	// 3. 锁定库存（乐观锁：stock > 0, status = 1）
	result := s.db.WithContext(ctx).Model(&itemModel.Product{}).
		Where("id = ? AND stock > 0 AND status = 1", req.ProductId).
		Updates(map[string]interface{}{
			"stock":  gorm.Expr("stock - 1"),
			"status": 2, // 2=已锁定
		})
	if result.RowsAffected == 0 {
		s.db.WithContext(ctx).Model(log).Update("status", "FAILED")
		return &item.TryLockStockResp{Success: false, BranchId: branchID, Msg: "库存不足"}, nil
	}

	// 4. 标记已锁定
	s.db.WithContext(ctx).Model(log).Update("status", "LOCKED")

	return &item.TryLockStockResp{Success: true, BranchId: branchID}, nil
}

// ConfirmDeductStock 确认扣减库存（TCC Confirm 阶段）
func (s *ItemService) ConfirmDeductStock(ctx context.Context, req *item.ConfirmDeductStockReq) (*item.ConfirmDeductStockResp, error) {
	var log itemModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ? AND action_type = 'LOCK_STOCK'", req.TxId, req.BranchId).First(&log).Error
	if err != nil {
		return &item.ConfirmDeductStockResp{Success: false, Msg: "事务日志不存在"}, nil
	}
	if log.Status == "CONFIRMED" {
		return &item.ConfirmDeductStockResp{Success: true, Msg: "幂等返回"}, nil
	}

	// 确认：状态改为已确认，商品状态改为已售
	s.db.WithContext(ctx).Model(&itemModel.Product{}).Where("id = ?", log.ProductID).Updates(map[string]interface{}{
		"status": 0, // 0=已售
	})
	s.db.WithContext(ctx).Model(&log).Update("status", "CONFIRMED")

	cache.GetRedis().Del(ctx, cache.BuildKey("product", log.ProductID))

	return &item.ConfirmDeductStockResp{Success: true}, nil
}

// CancelUnlockStock 解锁库存（TCC Cancel 阶段）
func (s *ItemService) CancelUnlockStock(ctx context.Context, req *item.CancelUnlockStockReq) (*item.CancelUnlockStockResp, error) {
	var log itemModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ? AND action_type = 'LOCK_STOCK'", req.TxId, req.BranchId).First(&log).Error
	if err != nil {
		return &item.CancelUnlockStockResp{Success: false, Msg: "事务日志不存在"}, nil
	}
	if log.Status == "CANCELLED" {
		return &item.CancelUnlockStockResp{Success: true, Msg: "幂等返回"}, nil
	}
	if log.Status == "CONFIRMED" {
		return &item.CancelUnlockStockResp{Success: false, Msg: "事务已确认，无法回滚"}, nil
	}

	// 回滚：恢复库存
	s.db.WithContext(ctx).Model(&itemModel.Product{}).Where("id = ?", log.ProductID).Updates(map[string]interface{}{
		"stock":  gorm.Expr("stock + 1"),
		"status": 1,
	})
	s.db.WithContext(ctx).Model(&log).Updates(map[string]interface{}{
		"status":     "CANCELLED",
		"action_type": "UNLOCK_STOCK",
	})

	cache.GetRedis().Del(ctx, cache.BuildKey("product", log.ProductID))

	return &item.CancelUnlockStockResp{Success: true}, nil
}

// GetProductDetail 查询商品详情
func (s *ItemService) GetProductDetail(ctx context.Context, productID uint64) (*item.GetProductDetailResp, error) {
	var p itemModel.Product
	err := s.db.WithContext(ctx).First(&p, productID).Error
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return &item.GetProductDetailResp{
		Id: uint64(p.ID), SellerId: uint64(p.SellerID),
		Title: p.Title, Desc: p.Desc, Category: p.Category,
		Tag: p.Tag, Price: p.Price, OriginalPrice: p.OriginalPrice,
		Images: p.Images, Campus: p.Campus, Status: int32(p.Status),
		Stock: int32(p.Stock), CreatedAt: p.CreatedAt.UnixMilli(),
	}, nil
}

// ListProducts 批量查询商品
func (s *ItemService) ListProducts(ctx context.Context, req *item.ListProductsReq) (*item.ListProductsResp, error) {
	offset := int((req.Page - 1) * req.PageSize)
	limit := int(req.PageSize)
	if req.PageSize <= 0 {
		limit = 20
	}

	var total int64
	var products []itemModel.Product
	query := s.db.WithContext(ctx).Model(&itemModel.Product{}).Where("status = 1")
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Campus != "" {
		query = query.Where("campus LIKE ?", "%"+req.Campus+"%")
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR `desc` LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.SellerId > 0 {
		query = query.Where("seller_id = ?", req.SellerId)
	}
	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&products)

	items := make([]*item.GetProductDetailResp, 0, len(products))
	for _, p := range products {
		items = append(items, &item.GetProductDetailResp{
			Id: uint64(p.ID), SellerId: uint64(p.SellerID),
			Title: p.Title, Desc: p.Desc, Category: p.Category,
			Price: p.Price, Images: p.Images, Campus: p.Campus,
			Status: int32(p.Status), CreatedAt: p.CreatedAt.UnixMilli(),
		})
	}
	return &item.ListProductsResp{Items: items, Total: total}, nil
}
