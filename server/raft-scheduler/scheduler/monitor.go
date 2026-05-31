package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"pxx/idl/gen/trade"
)

// Monitor 延迟调度监控器
type Monitor struct {
	raftNode  *raft.Raft
	fsm       raftFSM
	tradeAddr string // Trade 服务 gRPC 地址
	// gRPC 长连接复用
	tradeConn   *grpc.ClientConn
	tradeClient trade.TradeClient
}

type raftFSM interface {
	Heap() *TaskHeap
}

func NewMonitor(raftNode *raft.Raft, fsm raftFSM, tradeAddr string) *Monitor {
	m := &Monitor{
		raftNode:  raftNode,
		fsm:       fsm,
		tradeAddr: tradeAddr,
	}
	// 启动时建立长连接
	m.connectTrade()
	return m
}

func (m *Monitor) connectTrade() {
	// F-5 修复：关闭旧连接防止泄漏
	if m.tradeConn != nil {
		m.tradeConn.Close()
		m.tradeConn = nil
		m.tradeClient = nil
	}

	var err error
	m.tradeConn, err = grpc.Dial(m.tradeAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("[Scheduler] dial trade %s failed: %v", m.tradeAddr, err)
		m.tradeConn = nil
		return
	}
	m.tradeClient = trade.NewTradeClient(m.tradeConn)
	log.Printf("[Scheduler] connected to trade service at %s", m.tradeAddr)
}

// Start 启动调度循环（仅 Leader 执行）
func (m *Monitor) Start(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	log.Println("[Scheduler] monitor started")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Scheduler] monitor stopped")
			if m.tradeConn != nil {
				m.tradeConn.Close()
			}
			return
		case <-ticker.C:
			m.processDueTasks()
		}
	}
}

func (m *Monitor) processDueTasks() {
	if m.raftNode.State() != raft.Leader {
		return // 非 Leader 不执行调度
	}

	now := time.Now()
	// ★ 关键修复：PeekDue 只读查看到期任务，不修改 FSM 状态
	dueTasks := m.fsm.Heap().PeekDue(now)
	if len(dueTasks) == 0 {
		return
	}

	if m.tradeClient == nil {
		m.connectTrade()
		if m.tradeClient == nil {
			log.Printf("[Scheduler] cannot connect to trade, retry later")
			return
		}
	}

	for _, task := range dueTasks {
		// 二次确认：任务可能已被其他节点的 REMOVE 处理
		if !m.fsm.Heap().Exists(task.TaskID) {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := m.tradeClient.CancelOrder(ctx, &trade.CancelOrderReq{
			OrderId: task.OrderID,
			Reason:  task.Reason,
		})
		cancel()

		if err == nil && resp != nil && resp.Success {
			// 成功 → 通过 Raft Apply REMOVE
			cmd := Command{Type: "REMOVE", TaskID: task.TaskID}
			data, _ := json.Marshal(cmd)
			future := m.raftNode.Apply(data, 5*time.Second)
			if err := future.Error(); err != nil {
				log.Printf("[Raft] apply REMOVE failed: %v", err)
			}
			log.Printf("[Scheduler] cancelled order %d (task %s)", task.OrderID, task.TaskID)
		} else {
			errMsg := "nil resp"
			if err != nil {
				errMsg = err.Error()
			} else if resp != nil {
				errMsg = resp.Msg
			}
			log.Printf("[Scheduler] cancel order %d failed: %s, retry in 5s", task.OrderID, errMsg)
			// ★ 关键修复：通过 Raft Apply 提交重试任务，而非直接操作 FSM
			retryCmd := Command{
				Type:      "ADD",
				TaskID:    task.TaskID,
				OrderID:   task.OrderID,
				TriggerAt: now.Add(5 * time.Second),
				Reason:    task.Reason,
			}
			data, _ := json.Marshal(retryCmd)
			future := m.raftNode.Apply(data, 5*time.Second)
			if err := future.Error(); err != nil {
				log.Printf("[Raft] apply retry ADD failed: %v", err)
			}
		}
	}
}
