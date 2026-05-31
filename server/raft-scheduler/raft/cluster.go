package raft

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// Config Raft cluster config
type Config struct {
	NodeID        string
	RaftDir       string
	BindAddr      string
	JoinAddr      string
	AdvertiseAddr string
	Bootstrap     bool
}

// JoinRequest 节点加入集群请求
type JoinRequest struct {
	NodeID string `json:"node_id"`
	Addr   string `json:"addr"`
}

// JoinResponse 节点加入集群响应
type JoinResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

func NewRaftNode(cfg Config, fsm raft.FSM) (*raft.Raft, error) {
	// Raft Config
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(cfg.NodeID)
	config.SnapshotInterval = 60 * time.Second
	config.SnapshotThreshold = 100

	// Raft Store
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.RaftDir, "raft-log.bolt"))
	if err != nil {
		return nil, fmt.Errorf("bolt log store: %w", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.RaftDir, "raft-stable.bolt"))
	if err != nil {
		return nil, fmt.Errorf("bolt stable store: %w", err)
	}

	// Snapshot Store
	snapshotStore, err := raft.NewFileSnapshotStore(cfg.RaftDir, 3, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: %w", err)
	}

	// Transport
	advertiseAddr := cfg.AdvertiseAddr
	if advertiseAddr == "" {
		advertiseAddr = cfg.BindAddr
	}
	addr, err := net.ResolveTCPAddr("tcp", advertiseAddr)
	if err != nil {
		return nil, fmt.Errorf("resolve addr: %w", err)
	}
	transport, err := raft.NewTCPTransport(cfg.BindAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("transport: %w", err)
	}

	// Create Raft node
	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("new raft: %w", err)
	}

	// Bootstrap cluster
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

	// Start join listener on a derived port (Raft port + 10000)
	joinListenAddr := deriveJoinAddr(cfg.BindAddr)
	go startJoinListener(raftNode, joinListenAddr)

	log.Printf("[Raft] node %s started at %s (join: %s)", cfg.NodeID, cfg.BindAddr, joinListenAddr)
	return raftNode, nil
}

// deriveJoinAddr 从 Raft bind 地址推导 Join 端口（+10000 偏移）
// e.g. "0.0.0.0:7000" → "0.0.0.0:17000"
func deriveJoinAddr(bindAddr string) string {
	_, portStr, err := net.SplitHostPort(bindAddr)
	if err != nil {
		return bindAddr
	}
	// 解析端口并加偏移
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	// 取主机部分
	host, _, _ := net.SplitHostPort(bindAddr)
	return fmt.Sprintf("%s:%d", host, port+10000)
}

// deriveJoinAddrFromRaft 从 Raft 地址推导对端 Join 端口
// 用于 Join 请求：把 Raft 端口转为 Join 端口
func deriveJoinAddrFromRaft(raftAddr string) string {
	host, portStr, err := net.SplitHostPort(raftAddr)
	if err != nil {
		return raftAddr
	}
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	return fmt.Sprintf("%s:%d", host, port+10000)
}

// startJoinListener 在指定地址上监听 Join 请求
func startJoinListener(rn *raft.Raft, addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("[Raft] join listener failed on %s: %v", addr, err)
		return
	}
	defer ln.Close()
	log.Printf("[Raft] join listener started on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[Raft] join accept error: %v", err)
			continue
		}
		go handleJoinRequest(rn, conn)
	}
}

// handleJoinRequest 处理单个 Join 请求
func handleJoinRequest(rn *raft.Raft, conn net.Conn) {
	defer conn.Close()

	// 设置读超时
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	var req JoinRequest
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&req); err != nil {
		resp := JoinResponse{Success: false, Msg: fmt.Sprintf("decode error: %v", err)}
		json.NewEncoder(conn).Encode(resp)
		return
	}

	// 检查当前节点是否是 Leader
	if rn.State() != raft.Leader {
		resp := JoinResponse{Success: false, Msg: fmt.Sprintf("not leader, current leader: %s", rn.Leader())}
		json.NewEncoder(conn).Encode(resp)
		return
	}

	// 添加到集群
	addr, err := net.ResolveTCPAddr("tcp", req.Addr)
	if err != nil {
		resp := JoinResponse{Success: false, Msg: fmt.Sprintf("resolve addr error: %v", err)}
		json.NewEncoder(conn).Encode(resp)
		return
	}

	future := rn.AddVoter(raft.ServerID(req.NodeID), raft.ServerAddress(addr.String()), 0, 0)
	if err := future.Error(); err != nil {
		resp := JoinResponse{Success: false, Msg: fmt.Sprintf("add voter error: %v", err)}
		json.NewEncoder(conn).Encode(resp)
		return
	}

	resp := JoinResponse{Success: true, Msg: "joined"}
	json.NewEncoder(conn).Encode(resp)
	log.Printf("[Raft] node %s (%s) joined cluster", req.NodeID, req.Addr)
}

// JoinCluster 连接到 Leader 的 Join 端口，发送加入请求
func JoinCluster(joinAddr, nodeID, advertiseAddr string) error {
	// 将 Leader 的 Raft 地址转换为 Join 端口地址
	// joinAddr 格式可能是 "host:port"（Raft 端口）
	joinListenAddr := deriveJoinAddrFromRaft(joinAddr)

	conn, err := net.DialTimeout("tcp", joinListenAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("dial join listener %s: %w", joinListenAddr, err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// 获取本节点的 advertise 地址（用于其他节点连接自己）
	addrToAnnounce := advertiseAddr
	if addrToAnnounce == "" {
		addrToAnnounce = "127.0.0.1:7000" // default fallback
	}

	req := JoinRequest{
		NodeID: nodeID,
		Addr:   addrToAnnounce,
	}

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return fmt.Errorf("send join request: %w", err)
	}

	var resp JoinResponse
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return fmt.Errorf("read join response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("join rejected: %s", resp.Msg)
	}

	return nil
}