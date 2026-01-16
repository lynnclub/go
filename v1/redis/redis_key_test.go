package redis

import (
	"testing"
)

// TestKeyConstants 测试键常量定义
func TestKeyConstants(t *testing.T) {
	// 测试KeyBase常量
	if KeyBase != "general:" {
		t.Errorf("期望KeyBase为'general:'，实际为'%s'", KeyBase)
	}

	// 测试KeyLock常量
	expectedKeyLock := "general:lock:"
	if KeyLock != expectedKeyLock {
		t.Errorf("期望KeyLock为'%s'，实际为'%s'", expectedKeyLock, KeyLock)
	}
}

// TestKeyLockComposition 测试KeyLock组合
func TestKeyLockComposition(t *testing.T) {
	// 确保KeyLock是由KeyBase组合而成
	expectedKeyLock := KeyBase + "lock:"
	if KeyLock != expectedKeyLock {
		t.Errorf("KeyLock应该由KeyBase组合而成。期望'%s'，实际'%s'", expectedKeyLock, KeyLock)
	}
}

// TestKeyLockUsage 测试KeyLock在实际使用中的表现
func TestKeyLockUsage(t *testing.T) {
	// 测试拼接锁名称
	lockName := "test_resource"
	fullKey := KeyLock + lockName
	expected := "general:lock:test_resource"

	if fullKey != expected {
		t.Errorf("锁键拼接错误。期望'%s'，实际'%s'", expected, fullKey)
	}
}

// TestKeyBaseUsage 测试KeyBase在实际使用中的表现
func TestKeyBaseUsage(t *testing.T) {
	// 测试KeyBase可以用于创建其他键
	testKeys := []struct {
		suffix   string
		expected string
	}{
		{"cache:", "general:cache:"},
		{"session:", "general:session:"},
		{"queue:", "general:queue:"},
	}

	for _, tc := range testKeys {
		result := KeyBase + tc.suffix
		if result != tc.expected {
			t.Errorf("键拼接错误。后缀'%s'，期望'%s'，实际'%s'", tc.suffix, tc.expected, result)
		}
	}
}

// TestKeyFormat 测试键格式符合规范
func TestKeyFormat(t *testing.T) {
	// 测试键以冒号结尾，方便拼接
	if KeyBase[len(KeyBase)-1] != ':' {
		t.Error("KeyBase应该以冒号结尾")
	}

	if KeyLock[len(KeyLock)-1] != ':' {
		t.Error("KeyLock应该以冒号结尾")
	}
}

// TestKeyUniqueness 测试键的唯一性
func TestKeyUniqueness(t *testing.T) {
	// 确保不同的键前缀是唯一的
	if KeyBase == KeyLock {
		t.Error("KeyBase和KeyLock不应该相同")
	}
}
