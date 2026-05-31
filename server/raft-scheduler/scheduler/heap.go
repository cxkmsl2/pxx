package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

// TaskEntry 延迟任务条目
type TaskEntry struct {
	TaskID    string    `json:"task_id"`
	OrderID   uint64    `json:"order_id"`
	TriggerAt time.Time `json:"trigger_at"`
	Reason    string    `json:"reason"`
}

// TaskHeap 最小堆（按 TriggerAt 排序）—— 线程安全
type TaskHeap struct {
	mu    sync.RWMutex
	items []*TaskEntry
	index map[string]int // taskID -> index
}

func NewTaskHeap() *TaskHeap {
	h := &TaskHeap{
		items: make([]*TaskEntry, 0),
		index: make(map[string]int),
	}
	heap.Init(h)
	return h
}

// ---------- container/heap.Interface (内部使用，不加锁，调用者负责锁) ----------

func (h *TaskHeap) Len() int           { return len(h.items) }
func (h *TaskHeap) Less(i, j int) bool { return h.items[i].TriggerAt.Before(h.items[j].TriggerAt) }
func (h *TaskHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.index[h.items[i].TaskID] = i
	h.index[h.items[j].TaskID] = j
}

func (h *TaskHeap) Push(x interface{}) {
	item := x.(*TaskEntry)
	h.index[item.TaskID] = len(h.items)
	h.items = append(h.items, item)
}

func (h *TaskHeap) Pop() interface{} {
	old := h.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	delete(h.index, item.TaskID)
	h.items = old[:n-1]
	return item
}

// ---------- 公开方法（线程安全） ----------

// Peek 查看堆顶（只读）
func (h *TaskHeap) Peek() *TaskEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.items) == 0 {
		return nil
	}
	return h.items[0]
}

// Add 添加或更新任务（支持 Upsert）
func (h *TaskHeap) Add(entry *TaskEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if idx, exists := h.index[entry.TaskID]; exists {
		// 已存在：更新 TriggerAt 并重排堆
		h.items[idx].TriggerAt = entry.TriggerAt
		h.items[idx].Reason = entry.Reason
		heap.Fix(h, idx)
		return
	}
	heap.Push(h, entry)
}

// Remove 移除任务
func (h *TaskHeap) Remove(taskID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if idx, exists := h.index[taskID]; exists {
		heap.Remove(h, idx)
	}
}

// PopIfDue 弹出所有到期的任务
func (h *TaskHeap) PopIfDue(now time.Time) []*TaskEntry {
	h.mu.Lock()
	defer h.mu.Unlock()
	var due []*TaskEntry
	for {
		top := h.PeekUnsafe()
		if top == nil || top.TriggerAt.After(now) {
			break
		}
		due = append(due, heap.Pop(h).(*TaskEntry))
	}
	return due
}

// PeekDueUnsafe 只读：返回到期任务列表，不弹出（调用者需持有锁）
func (h *TaskHeap) PeekDueUnsafe(now time.Time) []*TaskEntry {
	var due []*TaskEntry
	for _, item := range h.items {
		if !item.TriggerAt.After(now) {
			due = append(due, item)
		}
	}
	return due
}

// PeekDue 只读：返回到期任务列表，不弹出（线程安全）
func (h *TaskHeap) PeekDue(now time.Time) []*TaskEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.PeekDueUnsafe(now)
}

// PeekUnsafe 查看堆顶（内部使用，不加锁）
func (h *TaskHeap) PeekUnsafe() *TaskEntry {
	if len(h.items) == 0 {
		return nil
	}
	return h.items[0]
}

// ListAll 返回所有任务（线程安全，用于快照）
func (h *TaskHeap) ListAll() []*TaskEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	entries := make([]*TaskEntry, len(h.items))
	// 深拷贝：复制 TaskEntry 值而非指针，避免快照时原始数据被修改
	for i, item := range h.items {
		copied := *item
		entries[i] = &copied
	}
	return entries
}

// Exists 检查任务是否存在
func (h *TaskHeap) Exists(taskID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.index[taskID]
	return exists
}
