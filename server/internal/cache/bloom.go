package cache

import (
	"hash/fnv"
	"sync"
)

type BloomFilter struct {
	bitset  []uint64
	size    uint32
	hashCnt uint32
	mu      sync.RWMutex
}

var bloom *BloomFilter

func InitBloomFilter(size uint32, hashCnt uint32) {
	bloom = &BloomFilter{
		bitset:  make([]uint64, (size+63)/64),
		size:    size,
		hashCnt: hashCnt,
	}
}

func GetBloomFilter() *BloomFilter {
	if bloom == nil {
		InitBloomFilter(1000000, 3)
	}
	return bloom
}

func (bf *BloomFilter) Add(key string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for i := uint32(0); i < bf.hashCnt; i++ {
		idx := bf.hash(key, i) % bf.size
		bf.bitset[idx/64] |= 1 << (idx % 64)
	}
}

func (bf *BloomFilter) MightExist(key string) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	for i := uint32(0); i < bf.hashCnt; i++ {
		idx := bf.hash(key, i) % bf.size
		if bf.bitset[idx/64]&(1<<(idx%64)) == 0 {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) hash(key string, seed uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte{byte(seed)})
	h.Write([]byte(key))
	return h.Sum32()
}
