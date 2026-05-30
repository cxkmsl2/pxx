package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"pxx/idl/gen/account"
	"pxx/idl/gen/item"
	"pxx/idl/gen/trade"
	tradeModel "pxx/trade/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TryResult Try 阶段的结果
type TryResult struct {
	name     string
	success  bool
	branchID uint64
	msg      string
}

type OrderService struct {
	db          *gorm.DB
	accountCli  account.AccountClient
	itemCli     item.ItemClient
}

func NewOrderService(db *gorm.DB, ac account.AccountClient, ic item.ItemClient) *OrderService {
	return &OrderService{
		db:         db,
		accountCli: ac,
		itemCli:    ic,
	}
}

// CreateOrder TCC 两阶段提交创建订单
// 1. 生成 TxID，写事务状态 TRYING
// 2. 并发调用 Account.TryFreeze + Item.TryLockStock
// 3. 都成功 → Confirm；任一失败 → Cancel
func (s *OrderService) CreateOrder(ctx context.Context, req *trade.CreateOrderReq) (*trade.CreateOrderResp, error) {
	txID := uuid.New().String()

	// 1. 写入事务状态 (TRYING)
	s.db.WithContext(ctx).Create(&tradeModel.TccTransaction{
		TxID:   txID,
		Status: "TRYING",
	})

	// 2. 并发 Try 阶段（带超时）
	tryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	results := make(chan TryResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	// 并发调用 Account.TryFreeze
	go func() {
		defer wg.Done()
		resp, err := s.accountCli.TryFreeze(tryCtx, &account.TryFreezeReq{
			TxId:   txID,
			UserId: req.BuyerId,
			Amount: req.Amount,
		})
		if err != nil || resp == nil || !resp.Success {
			msg := "冻结失败"
			if err != nil {
				msg = err.Error()
			} else if resp != nil {
				msg = resp.Msg
			}
			results <- TryResult{name: "account", success: false, msg: msg}
			return
		}
		results <- TryResult{name: "account", success: true, branchID: resp.BranchId}
	}()

	// 并发调用 Item.TryLockStock
	go func() {
		defer wg.Done()
		resp, err := s.itemCli.TryLockStock(tryCtx, &item.TryLockStockReq{
			TxId:      txID,
			ProductId: req.ProductId,
			UserId:    req.BuyerId,
		})
		if err != nil || resp == nil || !resp.Success {
			msg := "锁库存失败"
			if err != nil {
				msg = err.Error()
			} else if resp != nil {
				msg = resp.Msg
			}
			results <- TryResult{name: "item", success: false, msg: msg}
			return
		}
		results <- TryResult{name: "item", success: true, branchID: resp.BranchId}
	}()

	// 等待两个并发任务完成
	wg.Wait()
	close(results)

	// 3. 收集结果
	var accountResult, itemResult TryResult
	allSuccess := true
	for r := range results {
		switch r.name {
		case "account":
			accountResult = r
		case "item":
			itemResult = r
		}
		if !r.success {
			allSuccess = false
		}
	}

	// 4. 根据结果执行 Confirm 或 Cancel
	if allSuccess {
		return s.doConfirm(ctx, txID, req, accountResult.branchID, itemResult.branchID)
	} else {
		return s.doCancle(ctx, txID, accountResult, itemResult)
	}
}

// doConfirm 第二阶段：确认（所有 Try 成功）
func (s *OrderService) doConfirm(ctx context.Context, txID string, req *trade.CreateOrderReq,
	accountBranchID, itemBranchID uint64) (*trade.CreateOrderResp, error) {

	// 异步 Confirm
	var confirmWg sync.WaitGroup
	confirmWg.Add(2)

	go func() {
		defer confirmWg.Done()
		if resp, err := s.accountCli.ConfirmDeduct(context.Background(),
			&account.ConfirmDeductReq{TxId: txID, BranchId: accountBranchID}); err != nil {
			log.Printf("[TCC] account confirm failed: %v", err)
		} else {
			log.Printf("[TCC] account confirm: %v", resp.Success)
		}
	}()

	go func() {
		defer confirmWg.Done()
		if resp, err := s.itemCli.ConfirmDeductStock(context.Background(),
			&item.ConfirmDeductStockReq{TxId: txID, BranchId: itemBranchID}); err != nil {
			log.Printf("[TCC] item confirm failed: %v", err)
		} else {
			log.Printf("[TCC] item confirm: %v", resp.Success)
		}
	}()

	// 创建订单记录
	order := &tradeModel.Order{
		OrderNo:   fmt.Sprintf("PX%d%06d", time.Now().UnixMilli()%100000, req.ProductId),
		BuyerID:   uint(req.BuyerId),
		SellerID:  0, // 需要查询 product 来获取 seller_id
		ProductID: uint(req.ProductId),
		Amount:    req.Amount,
		Status:    tradeModel.OrderStatusPaid,
		TxID:      txID,
		PaidAt:    timePtr(time.Now()),
	}
	s.db.WithContext(ctx).Create(order)

	// 更新事务状态
	s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
		Where("tx_id = ?", txID).Update("status", "CONFIRMED")

	confirmWg.Wait()

	return &trade.CreateOrderResp{
		TxId:    txID,
		OrderId: uint64(order.ID),
		Status:  1, // SUCCESS
		Msg:     "下单成功",
	}, nil
}

// doCancle 第二阶段：取消（有 Try 失败）
func (s *OrderService) doCancle(ctx context.Context, txID string,
	accountResult, itemResult TryResult) (*trade.CreateOrderResp, error) {

	// 对成功的 Try 发起 Cancel
	if accountResult.success {
		go s.accountCli.CancelUnfreeze(context.Background(),
			&account.CancelUnfreezeReq{TxId: txID, BranchId: accountResult.branchID})
	}
	if itemResult.success {
		go s.itemCli.CancelUnlockStock(context.Background(),
			&item.CancelUnlockStockReq{TxId: txID, BranchId: itemResult.branchID})
	}

	// 更新事务状态
	s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
		Where("tx_id = ?", txID).Update("status", "CANCELLED")

	return &trade.CreateOrderResp{
		TxId:   txID,
		Status: 2, // CLOSED
		Msg:    "下单失败：" + accountResult.msg + " / " + itemResult.msg,
	}, nil
}

// GetOrder 查询订单详情
func (s *OrderService) GetOrder(ctx context.Context, orderID uint64) (*trade.GetOrderResp, error) {
	var order tradeModel.Order
	err := s.db.WithContext(ctx).First(&order, orderID).Error
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	resp := &trade.GetOrderResp{
		Id:        uint64(order.ID),
		OrderNo:   order.OrderNo,
		BuyerId:   uint64(order.BuyerID),
		SellerId:  uint64(order.SellerID),
		ProductId: uint64(order.ProductID),
		Amount:    order.Amount,
		Status:    int32(order.Status),
		TxId:      order.TxID,
		CreatedAt: order.CreatedAt.UnixMilli(),
	}
	if order.PaidAt != nil {
		resp.PaidAt = order.PaidAt.UnixMilli()
	}
	return resp, nil
}

// CancelOrder 取消订单（由 Raft 调度器调用）
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint64, reason string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&tradeModel.Order{}).Where("id = ? AND status = ?", orderID, tradeModel.OrderStatusPaid).
		Updates(map[string]interface{}{
			"status":      tradeModel.OrderStatusCancelled,
			"cancelled_at": &now,
		}).Error
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint64, status int32) error {
	return s.db.WithContext(ctx).Model(&tradeModel.Order{}).Where("id = ?", orderID).
		Update("status", status).Error
}

func timePtr(t time.Time) *time.Time {
	return &t
}
