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
}

type raftFSM interface {
	Heap() *TaskHeap
}

func NewMonitor(raftNode *raft.Raft, fsm raftFSM, tradeAddr string) *Monitor {
	return &Monitor{
		raftNode:  raftNode,
		fsm:       fsm,
		tradeAddr: tradeAddr,
	}
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
	dueTasks := m.fsm.Heap().PopIfDue(now)
	if len(dueTasks) == 0 {
		return
	}

	// 连接 Trade 服务并取消订单
	conn, err := grpc.Dial(m.tradeAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("[Scheduler] dial trade %s failed: %v", m.tradeAddr, err)
		// 失败时把任务重新放回去
		for _, t := range dueTasks {
			m.fsm.Heap().Add(t)
		}
		return
	}
	defer conn.Close()

	client := trade.NewTradeClient(conn)
	for _, task := range dueTasks {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := client.CancelOrder(ctx, &trade.CancelOrderReq{
			OrderId: task.OrderID,
			Reason:  task.Reason,
		})
		cancel()

		if err == nil && resp.Success {
			// 提交 REMOVE 日志到 Raft
			cmd := Command{Type: "REMOVE", TaskID: task.TaskID}
			data, _ := json.Marshal(cmd)
			future := m.raftNode.Apply(data, 5*time.Second)
			if err := future.Error(); err != nil {
				log.Printf("[Raft] apply REMOVE failed: %v", err)
			}
			log.Printf("[Scheduler] cancelled order %d (task %s)", task.OrderID, task.TaskID)
		} else {
			log.Printf("[Scheduler] cancel order %d failed: %v/%s", task.OrderID, err, resp)
			// 失败重试：放回堆中，5秒后重试
			retryTask := *task
			retryTask.TriggerAt = time.Now().Add(5 * time.Second)
			m.fsm.Heap().Add(&retryTask)
		}
	}
}
