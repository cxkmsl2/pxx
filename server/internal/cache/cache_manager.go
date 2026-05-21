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

type CacheManager struct {
	local  *LocalCache
	sg     singleflight.Group
}

func NewCacheManager() *CacheManager {
	return &CacheManager{local: GetLocalCache()}
}

// GetOrLoad 三层读取：local -> redis -> db
// 使用 singleflight 防止缓存击穿
func (cm *CacheManager) GetOrLoad(
	ctx context.Context,
	key string,
	dest interface{},
	redisTTL time.Duration,
	loader func(context.Context) (interface{}, error),
) error {
	// L2: 本地缓存
	if val, ok := cm.local.Get(key); ok {
		if s, isStr := val.(string); isStr && s == NullMarker {
			return ErrNotFound
		}
		b, _ := json.Marshal(val)
		return json.Unmarshal(b, dest)
	}

	// L3: Redis
	if err := GetJSON(ctx, key, dest); err == nil {
		cm.local.Set(key, dest, 5*time.Minute)
		return nil
	}

	// singleflight 防击穿
	v, err, _ := cm.sg.Do(key, func() (interface{}, error) {
		// 二次检查 Redis
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
			if rdb != nil { _ = rdb.Set(ctx, key, NullMarker, NullTTL).Err() }
			return NullMarker, nil
		}

		// 回写 Redis，TTL 加随机偏移防雪崩
		randTTL := redisTTL + time.Duration(rand.Int63n(int64(redisTTL/4)))
		if rdb != nil { _ = SetJSON(ctx, key, val, randTTL) }
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

// Invalidate 删除所有层级缓存
func (cm *CacheManager) Invalidate(ctx context.Context, keys ...string) {
	if cm.local == nil { return }
	for _, key := range keys {
		cm.local.Del(key)
		_ = DelKey(ctx, key)
	}
}

// BuildKey 统一缓存 key 构造
func BuildKey(entity string, id interface{}) string {
	return fmt.Sprintf("pxx:%s:%v", entity, id)
}
