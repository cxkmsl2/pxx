package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"
)

var (
	ErrNotFound    = errors.New("cache: not found")
	NullMarker     = "__NULL__"
	NullTTL        = 60 * time.Second
)

// CacheManager 双层缓存：Redis → singleflight → DB
// 已移除单机本地缓存层（local.go），仅依赖 Redis
type CacheManager struct {
	sg singleflight.Group
}

func NewCacheManager() *CacheManager {
	return &CacheManager{}
}

// GetOrLoad 两层读取：redis -> singleflight -> db
func (cm *CacheManager) GetOrLoad(
	ctx context.Context,
	key string,
	dest interface{},
	redisTTL time.Duration,
	loader func(context.Context) (interface{}, error),
) error {
	// L1: Redis
	if err := GetJSON(ctx, key, dest); err == nil {
		return nil
	}

	// singleflight 防击穿
	v, err, _ := cm.sg.Do(key, func() (interface{}, error) {
		// 二次检查 Redis（double check）
		var tmp interface{}
		if err2 := GetJSON(ctx, key, &tmp); err2 == nil {
			return tmp, nil
		}

		// 查 DB
		val, loadErr := loader(ctx)
		if loadErr != nil {
			return nil, loadErr
		}
		if val == nil {
			// 缓存空值防穿透
			if rdb != nil {
				_ = rdb.Set(ctx, key, NullMarker, NullTTL).Err()
			}
			return NullMarker, nil
		}

		// 回写 Redis，TTL 加随机偏移防雪崩
		randTTL := redisTTL + time.Duration(rand.Int63n(int64(redisTTL/4)))
		if rdb != nil {
			_ = SetJSON(ctx, key, val, randTTL)
		}
		return val, nil
	})

	if err != nil {
		return err
	}

	if s, ok := v.(string); ok && s == NullMarker {
		return ErrNotFound
	}

	b, _ := json.Marshal(v)
	return json.Unmarshal(b, dest)
}

// Invalidate 删除 Redis 缓存
func (cm *CacheManager) Invalidate(ctx context.Context, keys ...string) {
	for _, key := range keys {
		_ = DelKey(ctx, key)
	}
}

// BuildKey 统一缓存 key 构造
func BuildKey(entity string, id interface{}) string {
	return fmt.Sprintf("pxx:%s:%v", entity, id)
}
