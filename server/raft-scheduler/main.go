package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pxx/internal/mq"
	"pxx/pkg/config"
	raftNode "pxx/raft-scheduler/raft"
	"pxx/raft-scheduler/scheduler"

	"github.com/hashicorp/raft"
)

func main() {
	_ = config.Load()

	// 1. Raft 配置
	raftConfig := raftNode.Config{
		NodeID:        os.Getenv("RAFT_NODE_ID"),
		RaftDir:       os.Getenv("RAFT_DATA_DIR"),
		BindAddr:      os.Getenv("RAFT_BIND_ADDR"),
		JoinAddr:      os.Getenv("RAFT_JOIN_ADDR"),
		AdvertiseAddr: os.Getenv("RAFT_ADVERTISE_ADDR"),
		Bootstrap:     os.Getenv("RAFT_BOOTSTRAP") == "true",
	}
	if raftConfig.NodeID == "" {
		raftConfig.NodeID = "scheduler-1"
	}
	if raftConfig.RaftDir == "" {
		raftConfig.RaftDir = "/data/raft"
	}
	if raftConfig.BindAddr == "" {
		raftConfig.BindAddr = "0.0.0.0:7000"
	}

	fsm := raftNode.NewFSM()
	rn, err := raftNode.NewRaftNode(raftConfig, fsm)
	if err != nil {
		log.Fatalf("[Scheduler] raft init failed: %v", err)
	}

	// 2. 加入集群（非 Bootstrap 节点）
	if raftConfig.JoinAddr != "" && !raftConfig.Bootstrap {
		go func() {
			for i := 0; i < 30; i++ {
				log.Printf("[Scheduler] attempting to join cluster via %s (attempt %d/30)...",
					raftConfig.JoinAddr, i+1)
				if err := raftNode.JoinCluster(raftConfig.JoinAddr, raftConfig.NodeID,
					raftConfig.AdvertiseAddr); err != nil {
					log.Printf("[Scheduler] join failed: %v, retrying...", err)
					time.Sleep(2 * time.Second)
					continue
				}
				log.Printf("[Scheduler] successfully joined cluster as %s", raftConfig.NodeID)
				return
			}
			log.Fatalf("[Scheduler] FATAL: failed to join cluster after 30 attempts")
		}()
	}

	// 3. 启动调度监控
	tradeAddr := os.Getenv("TRADE_SERVICE_ADDR")
	if tradeAddr == "" {
		tradeAddr = "127.0.0.1:8080"
	}

	monitor := scheduler.NewMonitor(rn, fsm, tradeAddr)
	ctx, cancel := context.WithCancel(context.Background())
	go monitor.Start(ctx)

	// 4. 动态消费 Kafka 消息（仅 Leader 消费）
	kafkaBros := os.Getenv("KAFKA_BROKERS")
	if kafkaBros != "" {
		mq.InitKafka(kafkaBros)
		groupID := "raft-scheduler-shared-group"

		go func() {
			var consumerCtx context.Context
			var consumerCancel context.CancelFunc

			leaderCh := rn.LeaderCh()
			for {
				select {
				case isLeader := <-leaderCh:
					if isLeader {
						if consumerCancel != nil {
							consumerCancel()
						}
						log.Println("[Scheduler] became leader, starting Kafka consumer...")
						consumerCtx, consumerCancel = context.WithCancel(context.Background())
						go runLeaderConsumer(consumerCtx, kafkaBros, groupID, rn)
					} else {
						log.Println("[Scheduler] lost leadership, stopping Kafka consumer...")
						if consumerCancel != nil {
							consumerCancel()
							consumerCancel = nil
						}
					}
				case <-ctx.Done():
					if consumerCancel != nil {
						consumerCancel()
					}
					return
				}
			}
		}()
	}

	log.Println("[Scheduler] raft-scheduler started")

	// 5. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Scheduler] shutting down...")
	cancel()
	future := rn.Shutdown()
	if err := future.Error(); err != nil {
		log.Printf("[Scheduler] raft shutdown: %v", err)
	}
	log.Println("[Scheduler] shutdown complete")
}

// runLeaderConsumer Leader 节点运行的消费者
func runLeaderConsumer(ctx context.Context, brokers, groupID string, rn *raft.Raft) {
	mq.StartConsumer(ctx, brokers, mq.TopicOrderCreated, groupID, func(data []byte) error {
		var orderData map[string]interface{}
		if err := json.Unmarshal(data, &orderData); err != nil {
			return err
		}

		orderIDFloat, ok := orderData["ID"].(float64)
		if !ok {
			return nil
		}
		orderID := uint64(orderIDFloat)

		cmd := scheduler.Command{
			Type:      "ADD",
			TaskID:    fmt.Sprintf("timeout_%d", orderID),
			OrderID:   orderID,
			TriggerAt: time.Now().Add(15 * time.Minute),
			Reason:    "15分钟未支付自动取消",
		}
		cmdBytes, _ := json.Marshal(cmd)

		future := rn.Apply(cmdBytes, 5*time.Second)
		if err := future.Error(); err != nil {
			if rn.State() != raft.Leader {
				log.Printf("[Scheduler] lost leadership during apply, skipping order %d", orderID)
			} else {
				log.Printf("[Scheduler] raft add task failed (still leader): %v", err)
			}
			return nil
		}
		log.Printf("[Scheduler] add timeout task for order %d", orderID)
		return nil
	})

	<-ctx.Done()
	log.Println("[Scheduler] leader consumer stopped")
}