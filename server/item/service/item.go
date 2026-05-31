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
// 防悬挂与原子性保证
func (s *ItemService) TryLockStock(ctx context.Context, req *item.TryLockStockReq) (*item.TryLockStockResp, error) {
	branchID := branchID(req.TxId, "item")

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 幂等性与防悬挂检查
		var existing itemModel.TCCLog
		if err := tx.Where("tx_id = ? AND branch_id = ?", req.TxId, branchID).First(&existing).Error; err == nil {
			if existing.Status == "CANCELLED" {
				return fmt.Errorf("事务已取消")
			}
			return nil
		}

		// 2. 锁定库存（乐观锁）
		result := tx.Model(&itemModel.Product{}).
			Where("id = ? AND stock > 0 AND status = 1", req.ProductId).
			Updates(map[string]interface{}{
				"stock":  gorm.Expr("stock - 1"),
				"status": 2, // 2=已锁定
			})
		if result.RowsAffected == 0 {
			return fmt.Errorf("库存不足或商品状态不可用")
		}

		// 3. 插入 LOCKED 记录
		log := &itemModel.TCCLog{
			TxID:      req.TxId,
			BranchID:  branchID,
			ActionType: "LOCK_STOCK",
			Status:    "LOCKED",
			ProductID: req.ProductId,
			Payload:   "{}",
		}
		return tx.Create(log).Error
	})

	if err != nil {
		return &item.TryLockStockResp{Success: false, Msg: err.Error()}, nil
	}

	return &item.TryLockStockResp{Success: true, BranchId: branchID}, nil
}

// ConfirmDeductStock 确认扣减库存（TCC Confirm 阶段）
func (s *ItemService) ConfirmDeductStock(ctx context.Context, req *item.ConfirmDeductStockReq) (*item.ConfirmDeductStockResp, error) {
	var log itemModel.TCCLog
	err := s.db.WithContext(ctx).Where("tx_id = ? AND branch_id = ?", req.TxId, req.BranchId).First(&log).Error
	if err != nil {
		return &item.ConfirmDeductStockResp{Success: false, Msg: "事务日志不存在"}, nil
	}
	if log.Status == "CONFIRMED" {
		return &item.ConfirmDeductStockResp{Success: true, Msg: "幂等返回"}, nil
	}

	// 确认：状态改为已确认，商品状态改为已售
	s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Model(&itemModel.Product{}).Where("id = ? AND status = 2", log.ProductID).Updates(map[string]interface{}{
			"status": 0, // 0=已售
		})
		return tx.Model(&log).Update("status", "CONFIRMED").Error
	})

	cache.GetRedis().Del(ctx, cache.BuildKey("product", log.ProductID))
	return &item.ConfirmDeductStockResp{Success: true}, nil
}

// CancelUnlockStock 解锁库存（TCC Cancel 阶段）
func (s *ItemService) CancelUnlockStock(ctx context.Context, req *item.CancelUnlockStockReq) (*item.CancelUnlockStockResp, error) {
	bid := req.BranchId
	if bid == 0 {
		bid = branchID(req.TxId, "item")
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var log itemModel.TCCLog
		err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("tx_id = ? AND branch_id = ?", req.TxId, bid).First(&log).Error
		
		if err != nil {
			// 防空回滚：插入一条 CANCELLED 记录，防止后续延迟的 Try 成功
			emptyLog := &itemModel.TCCLog{
				TxID:       req.TxId,
				BranchID:   bid,
				ActionType: "UNLOCK_STOCK",
				Status:     "CANCELLED",
				Payload:    "{}",
			}
			return tx.Create(emptyLog).Error
		}

		if log.Status == "CANCELLED" {
			return nil
		}
		if log.Status == "CONFIRMED" {
			return fmt.Errorf("事务已确认，无法回滚")
		}

		// 只有 LOCKED 状态才需要回滚库存
		if log.Status == "LOCKED" {
			tx.Model(&itemModel.Product{}).Where("id = ? AND status = 2", log.ProductID).Updates(map[string]interface{}{
				"stock":  gorm.Expr("stock + 1"),
				"status": 1,
			})
		}

		err = tx.Model(&log).Updates(map[string]interface{}{
			"status":     "CANCELLED",
			"action_type": "UNLOCK_STOCK",
		}).Error
		if err == nil {
			// 在事务成功后删除缓存
			cache.GetRedis().Del(ctx, cache.BuildKey("product", log.ProductID))
		}
		return err
	})

	if err != nil {
		return &item.CancelUnlockStockResp{Success: false, Msg: err.Error()}, nil
	}

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
	// ★ 修复 Bug-11: 包含 status=2（已锁定）的商品，前端通过 status 字段区分展示
	query := s.db.WithContext(ctx).Model(&itemModel.Product{}).Where("status IN ?", []int{1, 2})
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