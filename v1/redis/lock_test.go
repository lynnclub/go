package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// TestLock 测试加锁
func TestLock(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// 配置Redis
	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default", option)

	lockName := "test_lock"

	// 第一次加锁应该成功
	if !Lock(lockName, 10) {
		t.Error("第一次加锁失败")
	}

	// 再次加锁应该失败（锁已存在）
	if Lock(lockName, 10) {
		t.Error("重复加锁应该失败")
	}

	// 解锁
	Unlock(lockName)

	// 解锁后再次加锁应该成功
	if !Lock(lockName, 10) {
		t.Error("解锁后加锁失败")
	}

	// 清理
	Unlock(lockName)
}

// TestLockExpire 测试锁过期
func TestLockExpire(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_expire", option)

	// 清理之前的连接池，强制重新连接
	pool.Delete("default")
	options["default"] = options["default_expire"]

	lockName := "test_lock_expire"

	// 加锁，1秒过期
	if !Lock(lockName, 1) {
		t.Error("加锁失败")
	}

	// 立即加锁应该失败
	if Lock(lockName, 1) {
		t.Error("锁未生效，重复加锁应该失败")
	}

	// 使用miniredis的FastForward来模拟时间流逝
	s.FastForward(2 * time.Second)

	// 过期后应该可以加锁
	if !Lock(lockName, 1) {
		t.Error("锁过期后加锁失败")
	}

	// 清理
	Unlock(lockName)
}

// TestUnlock 测试解锁
func TestUnlock(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_unlock", option)

	pool.Delete("default")
	options["default"] = options["default_unlock"]

	lockName := "test_unlock"

	// 加锁
	if !Lock(lockName, 10) {
		t.Error("加锁失败")
	}

	// 解锁
	Unlock(lockName)

	// 检查是否已解锁（应该可以再次加锁）
	if !Lock(lockName, 10) {
		t.Error("解锁后无法加锁")
	}

	// 清理
	Unlock(lockName)
}

// TestUnlockNonExistent 测试解锁不存在的锁
func TestUnlockNonExistent(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_unlock_ne", option)

	pool.Delete("default")
	options["default"] = options["default_unlock_ne"]

	// 解锁不存在的锁不应该引发错误
	Unlock("non_existent_lock")
}

// TestLockWithDifferentNames 测试不同名称的锁互不影响
func TestLockWithDifferentNames(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_different", option)

	pool.Delete("default")
	options["default"] = options["default_different"]

	lock1 := "lock_1"
	lock2 := "lock_2"

	// 同时锁定两个不同的名称
	if !Lock(lock1, 10) {
		t.Error("锁定lock_1失败")
	}
	if !Lock(lock2, 10) {
		t.Error("锁定lock_2失败")
	}

	// 两个锁都应该存在
	if Lock(lock1, 10) {
		t.Error("lock_1应该已被锁定")
	}
	if Lock(lock2, 10) {
		t.Error("lock_2应该已被锁定")
	}

	// 解锁lock_1不应该影响lock_2
	Unlock(lock1)
	if !Lock(lock1, 10) {
		t.Error("lock_1解锁后应该可以再次锁定")
	}
	if Lock(lock2, 10) {
		t.Error("lock_2应该仍然被锁定")
	}

	// 清理
	Unlock(lock1)
	Unlock(lock2)
}
