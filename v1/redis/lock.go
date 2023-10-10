package redis

import (
	"time"
)

// Lock 加锁，单位秒
func Lock(name string, expire time.Duration) bool {
	result, err := Use("").
		SetNX(Ctx, KeyLock+name, 1, expire*time.Second).
		Result()
	if err == nil && result {
		return true
	}

	return false
}

// Unlock 解锁
func Unlock(name string) {
	Use("").Del(Ctx, KeyLock+name)
}
