package raft

import (
	"encoding/json"

	"pxx/raft-scheduler/scheduler"

	"github.com/hashicorp/raft"
)

// Snapshot Raft 快照
type Snapshot struct {
	heap *scheduler.TaskHeap
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		entries := s.heap.ListAll()
		if err := json.NewEncoder(sink).Encode(entries); err != nil {
			return err
		}
		return sink.Close()
	}()
	if err != nil {
		sink.Cancel()
		return err
	}
	return nil
}

func (s *Snapshot) Release() {}
