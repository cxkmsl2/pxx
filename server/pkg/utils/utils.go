package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func RandomTTL(base time.Duration, jitter time.Duration) time.Duration {
	return base + time.Duration(rand.Int63n(int64(jitter)))
}

func Paginate(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset = (page - 1) * pageSize
	limit = pageSize
	return
}
