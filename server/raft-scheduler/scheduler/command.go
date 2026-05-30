package scheduler

import "time"

// Command Raft 日志命令
type Command struct {
	Type      string    `json:"type"`       // ADD / REMOVE
	TaskID    string    `json:"task_id"`
	OrderID   uint64    `json:"order_id"`
	TriggerAt time.Time `json:"trigger_at"`
	Reason    string    `json:"reason"`
}

// ListAll 返回所有任务
func (h *TaskHeap) ListAll() []*TaskEntry {
	entries := make([]*TaskEntry, len(h.items))
	copy(entries, h.items)
	return entries
}
