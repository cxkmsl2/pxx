package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pxx/pkg/config"
	raftNode "pxx/raft-scheduler/raft"
	"pxx/raft-scheduler/scheduler"
	"pxx/raft-scheduler/client"
)

func main() {
	_ = config.Load()

	// 1. Raft 配置
	raftConfig := raftNode.Config{
		NodeID:    os.Getenv("RAFT_NODE_ID"),
		RaftDir:   os.Getenv("RAFT_DATA_DIR"),
		BindAddr:  os.Getenv("RAFT_BIND_ADDR"),
		JoinAddr:  os.Getenv("RAFT_JOIN_ADDR"),
		Bootstrap: os.Getenv("RAFT_BOOTSTRAP") == "true",
	}
	if raftConfig.NodeID == "" { raftConfig.NodeID = "scheduler-1" }
	if raftConfig.RaftDir == "" { raftConfig.RaftDir = "/data/raft" }
	if raftConfig.BindAddr == "" { raftConfig.BindAddr = "0.0.0.0:7000" }

	fsm := raftNode.NewFSM()
	rn, err := raftNode.NewRaftNode(raftConfig, fsm)
	if err != nil {
		log.Fatalf("[Scheduler] raft init failed: %v", err)
	}

	// 2. 启动调度监控器
	tradeAddr := os.Getenv("TRADE_SERVICE_ADDR")
	if tradeAddr == "" { tradeAddr = "127.0.0.1:8080" }

	monitor := scheduler.NewMonitor(rn, fsm, tradeAddr)
	ctx, cancel := context.WithCancel(context.Background())
	go monitor.Start(ctx)

	// 3. 启动 gRPC 客户端用于接收外部添加任务请求
	tradeClient := client.NewTradeClient(tradeAddr)

	// 在订单创建时调用 Raft 添加延迟任务
	// 这是一个示例：生产环境应在 Trade 服务下单后直接调用
	_ = tradeClient

	log.Println("[Scheduler] raft-scheduler started")

	// 4. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Scheduler] shutting down...")
	cancel()
	future := rn.Shutdown()
	if err := future.Error(); err != nil {
		log.Printf("[Scheduler] raft shutdown: %v", err)
	}
}
