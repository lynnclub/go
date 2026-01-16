package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
)

// TestMaxMinGet 测试获取值
func TestMaxMinGet(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_get", option)

	// 清理并重置连接池
	pool.Delete("default")
	options["default"] = options["default_get"]

	m := &MaxMin{
		CacheKey: "test:",
		Name:     "get",
	}

	// 删除可能存在的旧数据
	m.Delete()

	// 初始值应该为0
	if val := m.Get(); val != 0 {
		t.Errorf("期望获取0，实际获取%d", val)
	}

	// 设置一个值
	err = Use("").Set(Ctx, m.CacheKey+m.Name, 100, 0).Err()
	if err != nil {
		t.Fatalf("设置值失败: %v", err)
	}

	// 获取值
	if val := m.Get(); val != 100 {
		t.Errorf("期望获取100，实际获取%d", val)
	}

	// 清理
	m.Delete()
}

// TestMaxMinSetMax 测试设置最大值
func TestMaxMinSetMax(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_max", option)

	pool.Delete("default")
	options["default"] = options["default_max"]

	m := &MaxMin{
		CacheKey: "test:",
		Name:     "set_max",
	}

	// 删除可能存在的旧数据
	m.Delete()

	// 第一次设置应该成功
	if err := m.SetMax(100); err != nil {
		t.Errorf("第一次设置最大值失败: %v", err)
	}

	if val := m.Get(); val != 100 {
		t.Errorf("期望获取100，实际获取%d", val)
	}

	// 设置更大的值应该成功
	if err := m.SetMax(200); err != nil {
		t.Errorf("设置更大的值失败: %v", err)
	}

	if val := m.Get(); val != 200 {
		t.Errorf("期望获取200，实际获取%d", val)
	}

	// 设置更小的值应该失败
	if err := m.SetMax(150); err == nil {
		t.Error("设置更小的值应该失败")
	}

	// 值应该保持不变
	if val := m.Get(); val != 200 {
		t.Errorf("期望获取200，实际获取%d", val)
	}

	// 设置相等的值应该失败
	if err := m.SetMax(200); err == nil {
		t.Error("设置相等的值应该失败")
	}

	// 清理
	m.Delete()
}

// TestMaxMinSetMin 测试设置最小值
func TestMaxMinSetMin(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_min", option)

	pool.Delete("default")
	options["default"] = options["default_min"]

	m := &MaxMin{
		CacheKey: "test:",
		Name:     "set_min",
	}

	// 删除可能存在的旧数据
	m.Delete()

	// 初始设置一个值
	Use("").Set(Ctx, m.CacheKey+m.Name, 100, 0)

	// 设置更小的值应该成功
	if err := m.SetMin(50); err != nil {
		t.Errorf("设置更小的值失败: %v", err)
	}

	if val := m.Get(); val != 50 {
		t.Errorf("期望获取50，实际获取%d", val)
	}

	// 设置更小的值应该成功
	if err := m.SetMin(25); err != nil {
		t.Errorf("设置更小的值失败: %v", err)
	}

	if val := m.Get(); val != 25 {
		t.Errorf("期望获取25，实际获取%d", val)
	}

	// 设置更大的值应该失败
	if err := m.SetMin(30); err == nil {
		t.Error("设置更大的值应该失败")
	}

	// 值应该保持不变
	if val := m.Get(); val != 25 {
		t.Errorf("期望获取25，实际获取%d", val)
	}

	// 设置相等的值应该失败
	if err := m.SetMin(25); err == nil {
		t.Error("设置相等的值应该失败")
	}

	// 清理
	m.Delete()
}

// TestMaxMinDelete 测试删除
func TestMaxMinDelete(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_del", option)

	pool.Delete("default")
	options["default"] = options["default_del"]

	m := &MaxMin{
		CacheKey: "test:",
		Name:     "delete",
	}

	// 设置一个值
	Use("").Set(Ctx, m.CacheKey+m.Name, 100, 0)

	if val := m.Get(); val != 100 {
		t.Errorf("期望获取100，实际获取%d", val)
	}

	// 删除
	m.Delete()

	// 删除后应该获取0
	if val := m.Get(); val != 0 {
		t.Errorf("删除后期望获取0，实际获取%d", val)
	}
}

// TestMaxMinEmptyName 测试空名称
func TestMaxMinEmptyName(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default_empty", option)

	pool.Delete("default")
	options["default"] = options["default_empty"]

	m := &MaxMin{
		CacheKey: "test:",
		Name:     "",
	}

	// 测试空名称也能正常工作
	if err := m.SetMax(100); err != nil {
		t.Errorf("设置失败: %v", err)
	}

	if val := m.Get(); val != 100 {
		t.Errorf("期望获取100，实际获取%d", val)
	}

	m.Delete()
}
