package cache

import (
	"sync"
	"time"
)

type localEntry struct {
	value      interface{}
	expireAt   time.Time
}

type LocalCache struct {
	mu    sync.RWMutex
	items map[string]*localEntry
}

var local *LocalCache

func InitLocalCache() {
	local = &LocalCache{
		items: make(map[string]*localEntry),
	}
	go local.cleanup()
}

func GetLocalCache() *LocalCache {
	return local
}

func (lc *LocalCache) Get(key string) (interface{}, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	entry, ok := lc.items[key]
	if !ok || time.Now().After(entry.expireAt) {
		return nil, false
	}
	return entry.value, true
}

func (lc *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.items[key] = &localEntry{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

func (lc *LocalCache) Del(key string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	delete(lc.items, key)
}

func (lc *LocalCache) cleanup() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		lc.mu.Lock()
		now := time.Now()
		for k, v := range lc.items {
			if now.After(v.expireAt) {
				delete(lc.items, k)
			}
		}
		lc.mu.Unlock()
	}
}
