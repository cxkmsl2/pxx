// Package aimd implements an Additive-Increase Multiplicative-Decrease
// adaptive concurrency limiter based on real-time RTT and error rate.
//
// AIMD is a congestion control algorithm originally from TCP. When the
// backend is healthy (low RTT, low errors), the limiter slowly increases
// the max concurrency window (+5 per epoch). When the backend shows
// congestion (high RTT or high error rate), the window is halved.
//
// This prevents cascading failures: under high load, the limiter sheds
// excess requests before they hit the database, keeping the system stable.
package aimd

import (
	"log"
	"math"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// Config defines the AIMD limiter parameters.
type Config struct {
	MinWindow        int64
	MaxWindow        int64
	InitialWindow    int64
	Epoch            time.Duration
	AddStep          int64
	RTTThresholdLow  time.Duration
	RTTThresholdHigh time.Duration
	ErrRateHigh      float64
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MinWindow:        20,
		MaxWindow:        1000,
		InitialWindow:    200,
		Epoch:            1 * time.Second,
		AddStep:          5,
		RTTThresholdLow:  200 * time.Millisecond,
		RTTThresholdHigh: 500 * time.Millisecond,
		ErrRateHigh:      0.05,
	}
}

type metricsBucket struct {
	totalCount int64
	errorCount int64
	rttSumNs   int64
}

type Limiter struct {
	cfg       Config
	MaxWindow int64
	inflight  atomic.Int64
	mu        sync.Mutex
	bucket    metricsBucket
	stopped   atomic.Bool
}

func New(cfg Config) *Limiter {
	l := &Limiter{cfg: cfg, MaxWindow: cfg.InitialWindow}
	go l.adjustLoop()
	return l
}

func (l *Limiter) Acquire() bool {
	if l.stopped.Load() {
		return true
	}
	current := l.inflight.Add(1)
	window := atomic.LoadInt64(&l.MaxWindow)
	if current > window {
		l.inflight.Add(-1)
		return false
	}
	return true
}

func (l *Limiter) Release(startTime time.Time, err error) {
	l.inflight.Add(-1)
	rtt := time.Since(startTime)
	l.mu.Lock()
	l.bucket.totalCount++
	l.bucket.rttSumNs += rtt.Nanoseconds()
	if err != nil {
		l.bucket.errorCount++
	}
	l.mu.Unlock()
}

func (l *Limiter) Stop() {
	l.stopped.Store(true)
}

func (l *Limiter) adjustLoop() {
	ticker := time.NewTicker(l.cfg.Epoch)
	defer ticker.Stop()
	for range ticker.C {
		if l.stopped.Load() {
			return
		}
		l.mu.Lock()
		bucket := l.bucket
		l.bucket = metricsBucket{}
		l.mu.Unlock()

		if bucket.totalCount == 0 {
			continue
		}

		avgRTT := time.Duration(bucket.rttSumNs / bucket.totalCount)
		errRate := float64(bucket.errorCount) / float64(bucket.totalCount)

		var window int64
		switch {
		case avgRTT < l.cfg.RTTThresholdLow && errRate < l.cfg.ErrRateHigh:
			window = atomic.AddInt64(&l.MaxWindow, l.cfg.AddStep)
			log.Printf("[aimd] ADD  window=%d avgRTT=%v errRate=%.1f%%", window, avgRTT, errRate*100)
		case avgRTT > l.cfg.RTTThresholdHigh || errRate > l.cfg.ErrRateHigh:
			for {
				current := atomic.LoadInt64(&l.MaxWindow)
				half := int64(math.Max(float64(l.cfg.MinWindow), float64(current)/2))
				if atomic.CompareAndSwapInt64(&l.MaxWindow, current, half) {
					window = half
					break
				}
			}
			log.Printf("[aimd] MULTI window=%d avgRTT=%v errRate=%.1f%% CONGESTION", window, avgRTT, errRate*100)
		default:
			window = atomic.LoadInt64(&l.MaxWindow)
		}

		if window > l.cfg.MaxWindow {
			atomic.StoreInt64(&l.MaxWindow, l.cfg.MaxWindow)
		}
		if window < l.cfg.MinWindow {
			atomic.StoreInt64(&l.MaxWindow, l.cfg.MinWindow)
		}
	}
}

// GinMiddleware returns a Gin handler that applies AIMD concurrency limiting.
func GinMiddleware(l *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !l.Acquire() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "服务器繁忙，请稍后再试",
			})
			c.Abort()
			return
		}
		start := time.Now()
		c.Next()
		var err error
		if c.Writer.Status() >= 500 {
			err = http.ErrAbortHandler
		}
		l.Release(start, err)
	}
}
