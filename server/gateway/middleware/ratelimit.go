package middleware

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// AIMDLimiter 基于 AIMD（和式增加/积式减少）的并发限流器
// 模拟 TCP 拥塞控制：健康时线性增加窗口，异常时指数回退
type AIMDLimiter struct {
	cwnd        atomic.Int64   // 拥塞窗口（当前允许的最大并发数）
	inflight    atomic.Int64   // 当前正在处理的请求数
	minCWND     int64          // 最小窗口
	maxCWND     int64          // 最大窗口
	rttAvg      atomic.Int64   // 平均 RTT（微秒）

	mu          sync.RWMutex
	lastErrTime time.Time
	errCount    int64

	// 健康探针列表
	probes []Probe
}

// Probe 健康探针
type Probe struct {
	Name       string
	CheckFunc  func(ctx context.Context) (rtt time.Duration, err error)
}

// ProbeResult 探针结果
type ProbeResult struct {
	Healthy bool
	RTT     time.Duration
}

// NewAIMDLimiter 创建 AIMD 限流器
func NewAIMDLimiter(min, max int64) *AIMDLimiter {
	l := &AIMDLimiter{
		minCWND: min,
		maxCWND: max,
	}
	l.cwnd.Store(min)
	return l
}

// AddProbe 添加健康探针
func (l *AIMDLimiter) AddProbe(name string, fn func(ctx context.Context) (time.Duration, error)) {
	l.mu.Lock()
	l.probes = append(l.probes, Probe{Name: name, CheckFunc: fn})
	l.mu.Unlock()
}

// Acquire 获取并发许可，返回 false 表示限流
func (l *AIMDLimiter) Acquire() bool {
	cwnd := l.cwnd.Load()
	inflight := l.inflight.Add(1)

	if inflight > cwnd {
		// 超过窗口，拒绝
		l.inflight.Add(-1)
		return false
	}
	return true
}

// Release 释放并发许可（请求结束时调用）
func (l *AIMDLimiter) Release() {
	l.inflight.Add(-1)
}

// runProbe 执行一轮健康探测并调整 CWND
func (l *AIMDLimiter) runProbe() {
	l.mu.RLock()
	probes := make([]Probe, len(l.probes))
	copy(probes, l.probes)
	l.mu.RUnlock()

	if len(probes) == 0 {
		return
	}

	healthy := true
	var totalRTT time.Duration

	for _, p := range probes {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		rtt, err := p.CheckFunc(ctx)
		if err != nil {
			log.Printf("[AIMD] probe %s failed: %v", p.Name, err)
			healthy = false
			break
		}
		totalRTT += rtt
	}

	avgRTT := totalRTT / time.Duration(len(probes))
	l.rttAvg.Store(int64(avgRTT))

	cwnd := l.cwnd.Load()
	if healthy {
		// 加性增（AI）：每轮 +10
		newCWND := cwnd + 10
		if newCWND > l.maxCWND {
			newCWND = l.maxCWND
		}
		l.cwnd.Store(newCWND)
		log.Printf("[AIMD] healthy, AI: %d -> %d (RTT=%v)", cwnd, newCWND, avgRTT)
	} else {
		// 乘性减（MD）：减半
		newCWND := cwnd / 2
		if newCWND < l.minCWND {
			newCWND = l.minCWND
		}
		l.cwnd.Store(newCWND)
		log.Printf("[AIMD] degraded, MD: %d -> %d (RTT=%v)", cwnd, newCWND, avgRTT)
		l.mu.Lock()
		l.lastErrTime = time.Now()
		l.errCount++
		l.mu.Unlock()
	}
}

// StartProbeLoop 启动定时探测循环（每秒一次）
func (l *AIMDLimiter) StartProbeLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.runProbe()
		}
	}
}

// Stats 返回限流器当前状态（用于监控）
func (l *AIMDLimiter) Stats() map[string]interface{} {
	inflight := l.inflight.Load()
	cwnd := l.cwnd.Load()
	return map[string]interface{}{
		"cwnd":         cwnd,
		"inflight":     inflight,
		"utilization":  float64(inflight) / float64(cwnd) * 100,
		"rtt_avg_us":   l.rttAvg.Load(),
	}
}

// AIMDRateLimit 返回 Gin 中间件
func AIMDRateLimit(limiter *AIMDLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Acquire() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "系统繁忙，请稍候",
			})
			return
		}
		defer limiter.Release()
		c.Next()
	}
}
