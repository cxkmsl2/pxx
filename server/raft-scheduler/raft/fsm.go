package raft

import (
	"encoding/json"
	"io"
	"log"

	"pxx/raft-scheduler/scheduler"

	"github.com/hashicorp/raft"
)

// FSM Raft 状态机
type FSM struct {
	heap *scheduler.TaskHeap
}

func NewFSM() *FSM {
	return &FSM{heap: scheduler.NewTaskHeap()}
}

func (f *FSM) Heap() *scheduler.TaskHeap {
	return f.heap
}

// Apply 应用 Raft 日志
func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	var cmd scheduler.Command
	if err := json.Unmarshal(logEntry.Data, &cmd); err != nil {
		log.Printf("[Raft] unmarshal failed: %v", err)
		return nil
	}

	switch cmd.Type {
	case "ADD":
		entry := &scheduler.TaskEntry{
			TaskID:    cmd.TaskID,
			OrderID:   cmd.OrderID,
			TriggerAt: cmd.TriggerAt,
			Reason:    cmd.Reason,
		}
		f.heap.Add(entry)
		log.Printf("[Raft] ADD %s order=%d at %v", cmd.TaskID, cmd.OrderID, cmd.TriggerAt)

	case "REMOVE":
		f.heap.Remove(cmd.TaskID)
		log.Printf("[Raft] REMOVE %s", cmd.TaskID)
	}

	return nil
}

// Snapshot 快照
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{heap: f.heap}, nil
}

// Restore 从快照恢复
func (f *FSM) Restore(rc io.ReadCloser) error {
	var entries []*scheduler.TaskEntry
	if err := json.NewDecoder(rc).Decode(&entries); err != nil {
		return err
	}
	f.heap = scheduler.NewTaskHeap()
	for _, e := range entries {
		f.heap.Add(e)
	}
	log.Printf("[Raft] restored %d tasks", len(entries))
	return rc.Close()
}
