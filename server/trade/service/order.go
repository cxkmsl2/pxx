package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"pxx/idl/gen/account"
	"pxx/idl/gen/item"
	"pxx/idl/gen/trade"
	tradeModel "pxx/trade/model"
	"pxx/internal/mq"

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

func generateBranchID(txID string, service string) uint64 {
	var hash uint64 = 5381
	for _, c := range txID + service {
		hash = ((hash << 5) + hash) + uint64(c)
	}
	return hash
}

func timePtr(t time.Time) *time.Time { return &t }

// CreateOrder TCC 两阶段提交创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *trade.CreateOrderReq) (*trade.CreateOrderResp, error) {
	txID := uuid.New().String()

	payload, _ := json.Marshal(req)
	if err := s.db.WithContext(ctx).Create(&tradeModel.TccTransaction{
		TxID:    txID,
		Status:  "TRYING",
		Payload: string(payload),
	}).Error; err != nil {
		log.Printf("[TCC] create tcc_transaction failed: %v", err)
	}

	tryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	results := make(chan TryResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		resp, err := s.accountCli.TryFreeze(tryCtx, &account.TryFreezeReq{
			TxId: txID, UserId: req.BuyerId, Amount: req.Amount,
		})
		if err != nil || resp == nil || !resp.Success {
			msg := "冻结失败"
			if err != nil { msg = err.Error() } else if resp != nil { msg = resp.Msg }
			results <- TryResult{name: "account", success: false, msg: msg}
			return
		}
		results <- TryResult{name: "account", success: true, branchID: resp.BranchId}
	}()

	go func() {
		defer wg.Done()
		resp, err := s.itemCli.TryLockStock(tryCtx, &item.TryLockStockReq{
			TxId: txID, ProductId: req.ProductId, UserId: req.BuyerId,
		})
		if err != nil || resp == nil || !resp.Success {
			msg := "锁库存失败"
			if err != nil { msg = err.Error() } else if resp != nil { msg = resp.Msg }
			results <- TryResult{name: "item", success: false, msg: msg}
			return
		}
		results <- TryResult{name: "item", success: true, branchID: resp.BranchId}
	}()

	wg.Wait()
	close(results)

	var accountResult, itemResult TryResult
	allSuccess := true
	for r := range results {
		switch r.name {
		case "account": accountResult = r
		case "item":    itemResult = r
		}
		if !r.success { allSuccess = false }
	}

	if allSuccess {
		s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
			Where("tx_id = ?", txID).Update("status", "CONFIRMING")
		return s.doConfirm(ctx, txID, req, accountResult.branchID, itemResult.branchID)
	}

	s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
		Where("tx_id = ?", txID).Update("status", "CANCELLING")
	return s.doCancle(ctx, txID, accountResult, itemResult)
}

// doConfirm 第二阶段：确认
func (s *OrderService) doConfirm(ctx context.Context, txID string, req *trade.CreateOrderReq,
	accountBranchID, itemBranchID uint64) (*trade.CreateOrderResp, error) {

	productResp, err := s.itemCli.GetProductDetail(ctx, &item.GetProductDetailReq{ProductId: req.ProductId})
	var sellerID uint64
	if err == nil && productResp != nil {
		sellerID = productResp.SellerId
	}

	var confirmWg sync.WaitGroup
	confirmWg.Add(2)
	var accountOK, itemOK bool

	go func() {
		defer confirmWg.Done()
		confirmCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := s.accountCli.ConfirmDeduct(confirmCtx, &account.ConfirmDeductReq{TxId: txID, BranchId: accountBranchID})
		accountOK = (err == nil && resp != nil && resp.Success)
		if !accountOK { log.Printf("[TCC] account confirm failed: err=%v", err) }
	}()

	go func() {
		defer confirmWg.Done()
		confirmCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := s.itemCli.ConfirmDeductStock(confirmCtx, &item.ConfirmDeductStockReq{TxId: txID, BranchId: itemBranchID})
		itemOK = (err == nil && resp != nil && resp.Success)
		if !itemOK { log.Printf("[TCC] item confirm failed: err=%v", err) }
	}()

	confirmWg.Wait()

	order := &tradeModel.Order{
		OrderNo:   fmt.Sprintf("PX%d%06d", time.Now().UnixMilli()%100000, req.ProductId),
		BuyerID:   uint(req.BuyerId), SellerID: uint(sellerID), ProductID: uint(req.ProductId),
		Amount:    req.Amount, Status: tradeModel.OrderStatusPaid,
		TxID:      txID, PaidAt: timePtr(time.Now()),
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err == nil {
		_ = mq.Publish(context.Background(), mq.TopicOrderCreated, order)
	} else {
		log.Printf("[TCC] create order idempotency check: %v", err)
	}

	if accountOK && itemOK {
		s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
			Where("tx_id = ?", txID).Update("status", "CONFIRMED")
		return &trade.CreateOrderResp{TxId: txID, OrderId: uint64(order.ID), Status: 1, Msg: "下单成功"}, nil
	}

	// M1 修复：Confirm RPC 失败时不返回 success，返回 processing 状态告知前端等待
	return &trade.CreateOrderResp{TxId: txID, OrderId: uint64(order.ID), Status: 3, Msg: "订单处理中，请稍后查看"}, nil
}

// doCancle 第二阶段：取消
func (s *OrderService) doCancle(ctx context.Context, txID string,
	accountResult, itemResult TryResult) (*trade.CreateOrderResp, error) {

	accountBranch := accountResult.branchID
	if accountBranch == 0 { accountBranch = generateBranchID(txID, "account") }
	itemBranch := itemResult.branchID
	if itemBranch == 0 { itemBranch = generateBranchID(txID, "item") }

	var cancelWg sync.WaitGroup
	cancelWg.Add(2)

	go func() {
		defer cancelWg.Done()
		cancelCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.accountCli.CancelUnfreeze(cancelCtx, &account.CancelUnfreezeReq{TxId: txID, BranchId: accountBranch})
	}()

	go func() {
		defer cancelWg.Done()
		cancelCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.itemCli.CancelUnlockStock(cancelCtx, &item.CancelUnlockStockReq{TxId: txID, BranchId: itemBranch})
	}()
	cancelWg.Wait()

	s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
		Where("tx_id = ?", txID).Update("status", "CANCELLED")

	return &trade.CreateOrderResp{
		TxId: txID, Status: 2,
		Msg: "下单失败：" + accountResult.msg + " / " + itemResult.msg,
	}, nil
}

// GetOrder 查询订单详情
func (s *OrderService) GetOrder(ctx context.Context, orderID uint64) (*trade.GetOrderResp, error) {
	var order tradeModel.Order
	if err := s.db.WithContext(ctx).First(&order, orderID).Error; err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	resp := &trade.GetOrderResp{
		Id: uint64(order.ID), OrderNo: order.OrderNo,
		BuyerId: uint64(order.BuyerID), SellerId: uint64(order.SellerID),
		ProductId: uint64(order.ProductID), Amount: order.Amount,
		Status: int32(order.Status), TxId: order.TxID,
		CreatedAt: order.CreatedAt.UnixMilli(),
	}
	if order.PaidAt != nil { resp.PaidAt = order.PaidAt.UnixMilli() }
	return resp, nil
}

// ListOrders 查询订单列表
func (s *OrderService) ListOrders(ctx context.Context, userID uint64, page, pageSize int) ([]*trade.GetOrderResp, int64, error) {
	var orders []tradeModel.Order
	var total int64
	db := s.db.WithContext(ctx).Model(&tradeModel.Order{}).Where("buyer_id = ?", userID)
	db.Count(&total)
	if err := db.Order("id desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	var items []*trade.GetOrderResp
	for _, o := range orders {
		r := &trade.GetOrderResp{
			Id: uint64(o.ID), OrderNo: o.OrderNo,
			BuyerId: uint64(o.BuyerID), SellerId: uint64(o.SellerID),
			ProductId: uint64(o.ProductID), Amount: o.Amount,
			Status: int32(o.Status), TxId: o.TxID,
			CreatedAt: o.CreatedAt.UnixMilli(),
		}
		if o.PaidAt != nil { r.PaidAt = o.PaidAt.UnixMilli() }
		items = append(items, r)
	}
	return items, total, nil
}

// CancelOrder 取消订单（由 Raft 调度器调用）
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint64, reason string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&tradeModel.Order{}).
		Where("id = ? AND status IN ?", orderID, []int8{tradeModel.OrderStatusPending, tradeModel.OrderStatusPaid}).
		Updates(map[string]interface{}{
			"status":       tradeModel.OrderStatusCancelled,
			"cancelled_at": &now,
		}).Error
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint64, status int32) error {
	return s.db.WithContext(ctx).Model(&tradeModel.Order{}).
		Where("id = ?", orderID).Update("status", status).Error
}

// ============================ TCC 悬挂事务看门狗 ============================

func (s *OrderService) StartTCCWatchdog(ctx context.Context) {
	const scanInterval = 30 * time.Second
	const hangingTimeout = 5 * time.Minute

	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupHangingTransactions(ctx, hangingTimeout)
		}
	}
}

func (s *OrderService) cleanupHangingTransactions(ctx context.Context, timeout time.Duration) {
	cutoff := time.Now().Add(-timeout)
	var hanging []tradeModel.TccTransaction
	s.db.WithContext(ctx).
		Where("status IN ? AND created_at < ?", []string{"TRYING", "CONFIRMING", "CANCELLING"}, cutoff).
		Find(&hanging)

	for _, tx := range hanging {
		accountBranch := generateBranchID(tx.TxID, "account")
		itemBranch := generateBranchID(tx.TxID, "item")

		var req trade.CreateOrderReq
		_ = json.Unmarshal([]byte(tx.Payload), &req)

		switch tx.Status {
		case "TRYING", "CANCELLING":
			go safeGo(func() {
				s.doCancle(context.Background(), tx.TxID,
					TryResult{name: "account", branchID: accountBranch},
					TryResult{name: "item", branchID: itemBranch})
			})
		case "CONFIRMING":
			go safeGo(func() {
				s.retryConfirm(context.Background(), tx.TxID, &req, accountBranch, itemBranch)
			})
		}
	}
}

// retryConfirm 看门狗专用：恢复 Confirm 调用并确保订单记录存在
func (s *OrderService) retryConfirm(ctx context.Context, txID string, req *trade.CreateOrderReq, accountBranchID, itemBranchID uint64) {
	log.Printf("[TCC-Watchdog] retrying confirm for tx %s", txID)

	var wg sync.WaitGroup
	wg.Add(2)
	var accountOK, itemOK bool

	go func() {
		defer wg.Done()
		confirmCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		resp, err := s.accountCli.ConfirmDeduct(confirmCtx, &account.ConfirmDeductReq{TxId: txID, BranchId: accountBranchID})
		accountOK = (err == nil && resp != nil && resp.Success)
	}()

	go func() {
		defer wg.Done()
		confirmCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		resp, err := s.itemCli.ConfirmDeductStock(confirmCtx, &item.ConfirmDeductStockReq{TxId: txID, BranchId: itemBranchID})
		itemOK = (err == nil && resp != nil && resp.Success)
	}()

	wg.Wait()

	productResp, _ := s.itemCli.GetProductDetail(ctx, &item.GetProductDetailReq{ProductId: req.ProductId})
	sellerID := uint64(0)
	if productResp != nil { sellerID = productResp.SellerId }

	order := &tradeModel.Order{
		OrderNo:   fmt.Sprintf("PX%d%06d", time.Now().UnixMilli()%100000, req.ProductId),
		BuyerID:   uint(req.BuyerId), SellerID: uint(sellerID), ProductID: uint(req.ProductId),
		Amount:    req.Amount, Status: tradeModel.OrderStatusPaid, TxID: txID, PaidAt: timePtr(time.Now()),
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err == nil {
		_ = mq.Publish(context.Background(), mq.TopicOrderCreated, order)
	} else {
		log.Printf("[TCC-Watchdog] order already exists or create failed: %v", err)
	}

	if accountOK && itemOK {
		s.db.WithContext(ctx).Model(&tradeModel.TccTransaction{}).
			Where("tx_id = ?", txID).Update("status", "CONFIRMED")
		log.Printf("[TCC-Watchdog] tx %s confirm recovered", txID)
	}
}

// safeGo 启动 goroutine 并捕获 panic，防止单个恢复任务崩溃看门狗
func safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[TCC-Watchdog] panic recovered: %v", r)
			}
		}()
		fn()
	}()
}
