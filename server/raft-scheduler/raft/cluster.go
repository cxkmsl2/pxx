package raft

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// Config Raft 集群配置
type Config struct {
	NodeID      string
	RaftDir     string
	BindAddr    string
	JoinAddr    string // 加入现有集群的地址
	Bootstrap   bool   // 是否初始化集群
}

func NewRaftNode(cfg Config, fsm raft.FSM) (*raft.Raft, error) {
	// Raft 配置
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(cfg.NodeID)
	config.SnapshotInterval = 60 * time.Second
	config.SnapshotThreshold = 100

	// Raft 存储
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.RaftDir, "raft-log.bolt"))
	if err != nil {
		return nil, fmt.Errorf("bolt log store: %w", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.RaftDir, "raft-stable.bolt"))
	if err != nil {
		return nil, fmt.Errorf("bolt stable store: %w", err)
	}

	// 快照存储
	snapshotStore, err := raft.NewFileSnapshotStore(cfg.RaftDir, 3, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: %w", err)
	}

	// 网络传输
	addr, err := net.ResolveTCPAddr("tcp", cfg.BindAddr)
	if err != nil {
		return nil, fmt.Errorf("resolve addr: %w", err)
	}
	transport, err := raft.NewTCPTransport(cfg.BindAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("transport: %w", err)
	}

	// 创建 Raft 节点
	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("new raft: %w", err)
	}

	// 是否引导集群
	if cfg.Bootstrap {
		cluster := raft.Configuration{
			Servers: []raft.Server{
				{ID: config.LocalID, Address: transport.LocalAddr()},
			},
		}
		future := raftNode.BootstrapCluster(cluster)
		if err := future.Error(); err != nil {
			log.Printf("[Raft] bootstrap warning: %v", err)
		}
	}

	// Note: 加入集群使用 AddVoter，在 Leader 节点上执行
	// 此处简化，生产环境需实现 discovery 机制
	log.Printf("[Raft] node %s started at %s", cfg.NodeID, cfg.BindAddr)
	return raftNode, nil
}
