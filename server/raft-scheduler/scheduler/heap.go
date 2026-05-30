package scheduler

import (
	"container/heap"
	"time"
)

// TaskEntry 延迟任务条目
type TaskEntry struct {
	TaskID    string    `json:"task_id"`
	OrderID   uint64    `json:"order_id"`
	TriggerAt time.Time `json:"trigger_at"`
	Reason    string    `json:"reason"`
}

// TaskHeap 最小堆（按 TriggerAt 排序）
type TaskHeap struct {
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

// Peek 查看堆顶
func (h *TaskHeap) Peek() *TaskEntry {
	if len(h.items) == 0 {
		return nil
	}
	return h.items[0]
}

// Add 添加任务
func (h *TaskHeap) Add(entry *TaskEntry) {
	if _, exists := h.index[entry.TaskID]; exists {
		return // 已存在
	}
	heap.Push(h, entry)
}

// Remove 移除任务
func (h *TaskHeap) Remove(taskID string) {
	if idx, exists := h.index[taskID]; exists {
		heap.Remove(h, idx)
	}
}

// PopIfDue 弹出所有到期的任务
func (h *TaskHeap) PopIfDue(now time.Time) []*TaskEntry {
	var due []*TaskEntry
	for {
		top := h.Peek()
		if top == nil || top.TriggerAt.After(now) {
			break
		}
		due = append(due, heap.Pop(h).(*TaskEntry))
	}
	return due
}
