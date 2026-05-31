package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"pxx/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// TokenBucket 简易滑动窗口限流中间件
func TokenBucket(rate int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate_limit:" + c.ClientIP()
		rdb := cache.GetRedis()
		ctx := c.Request.Context()

		now := time.Now().UnixMilli()
		window := int64(60 * 1000) // 1分钟窗口
		
		pipe := rdb.Pipeline()
		// 移除窗口外的记录
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-window))
		// 统计窗口内请求数
		pipe.ZCard(ctx, key)
		// 添加当前请求 (使用纳秒时间戳保证唯一性)
		pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", time.Now().UnixNano())})
		// 设置过期时间
		pipe.Expire(ctx, key, time.Minute)
		
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			log.Println("[RateLimit] Redis unavailable, rate limit bypassed")
			c.Next()
			return
		}

		count, _ := cmds[1].(*redis.IntCmd).Result()
		if int(count) > burst {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"code": 429, "msg": "请求过于频繁，请稍后再试"})
			return
		}
		c.Next()
	}
}