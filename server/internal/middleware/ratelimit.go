package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"pxx/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// TokenBucket 简易令牌桶限流中间件
func TokenBucket(rate int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate_limit:" + c.ClientIP()
		rdb := cache.GetRedis()
		ctx := context.Background()

		now := time.Now().Unix()
		pipe := rdb.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(burst)))
		pipe.ZCard(ctx, key)
		pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-%d", now, rand.Int())})
		pipe.Expire(ctx, key, time.Duration(burst)*time.Second)
		cmds, err := pipe.Exec(ctx)

		if err != nil {
			c.Next()
			return
		}

		count, _ := cmds[1].(*redis.IntCmd).Result()
		if int(count) > burst {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"code": 429, "msg": "请求过于频繁"})
			return
		}
		c.Next()
	}
}
