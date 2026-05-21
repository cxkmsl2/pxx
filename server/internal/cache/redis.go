package cache

import (
	"fmt"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis(addr, password string) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return rdb.Ping(ctx).Err()
}

func GetRedis() *redis.Client {
	return rdb
}

// SetJSON 序列化后写入 Redis，TTL 带随机偏移防雪崩
func SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if rdb == nil { return nil }
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, data, ttl).Err()
}

// GetJSON 从 Redis 读取并反序列化
func GetJSON(ctx context.Context, key string, dest interface{}) error {
	if rdb == nil { return fmt.Errorf("redis not connected") }
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// SetNXJSON 原子性写入（分布式锁）
func SetNXJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return rdb.SetNX(ctx, key, data, ttl).Result()
}

// DelKey 删除缓存
func DelKey(ctx context.Context, key string) error {
	if rdb == nil { return nil }
	return rdb.Del(ctx, key).Err()
}

// IncrBy 原子自增
func IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	return rdb.IncrBy(ctx, key, delta).Result()
}

// Lock 简易分布式锁
func Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return rdb.SetNX(ctx, "lock:"+key, 1, ttl).Result()
}

func Unlock(ctx context.Context, key string) error {
	return rdb.Del(ctx, "lock:"+key).Err()
}
