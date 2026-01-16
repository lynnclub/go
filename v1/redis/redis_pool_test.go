package redis

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// setupMiniRedis 创建一个miniredis实例用于测试
func setupMiniRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("无法启动miniredis: %v", err)
	}
	return s
}

// TestAdd 测试添加配置
func TestAdd(t *testing.T) {
	option := Option{
		Address:  []string{"localhost:6379"},
		Password: "",
		DB:       0,
	}
	Add("test_add", option)

	if _, ok := options["test_add"]; !ok {
		t.Error("添加配置失败")
	}

	// 验证默认值设置
	opt := options["test_add"]
	if opt.PoolSize != 100 {
		t.Errorf("期望PoolSize为100，实际为%d", opt.PoolSize)
	}
	if opt.MasterName != "mymaster" {
		t.Errorf("期望MasterName为mymaster，实际为%s", opt.MasterName)
	}
	if opt.ConnMaxIdleTime != 5*time.Minute {
		t.Errorf("期望ConnMaxIdleTime为5分钟，实际为%v", opt.ConnMaxIdleTime)
	}
}

// TestAddMap 测试从map添加配置
func TestAddMap(t *testing.T) {
	setting := map[string]interface{}{
		"address":            []interface{}{"localhost:6379"},
		"password":           "test_password",
		"db":                 1,
		"pool_size":          50,
		"min_idle_conns":     5,
		"max_idle_conns":     20,
		"conn_max_idle_time": "10m",
		"master_name":        "custom_master",
		"tls":                true,
	}
	AddMap("test_add_map", setting)

	opt, ok := options["test_add_map"]
	if !ok {
		t.Error("从map添加配置失败")
	}

	// 验证所有字段
	if opt.Password != "test_password" {
		t.Errorf("期望Password为test_password，实际为%s", opt.Password)
	}
	if opt.DB != 1 {
		t.Errorf("期望DB为1，实际为%d", opt.DB)
	}
	if opt.PoolSize != 50 {
		t.Errorf("期望PoolSize为50，实际为%d", opt.PoolSize)
	}
	if opt.MinIdleConns != 5 {
		t.Errorf("期望MinIdleConns为5，实际为%d", opt.MinIdleConns)
	}
	if opt.MaxIdleConns != 20 {
		t.Errorf("期望MaxIdleConns为20，实际为%d", opt.MaxIdleConns)
	}
	if opt.ConnMaxIdleTime != 10*time.Minute {
		t.Errorf("期望ConnMaxIdleTime为10分钟，实际为%v", opt.ConnMaxIdleTime)
	}
	if opt.MasterName != "custom_master" {
		t.Errorf("期望MasterName为custom_master，实际为%s", opt.MasterName)
	}
	if !opt.TLS {
		t.Error("期望TLS为true，实际为false")
	}
}

// TestAddMapWithIntDuration 测试从map添加配置（整数形式的duration）
func TestAddMapWithIntDuration(t *testing.T) {
	setting := map[string]interface{}{
		"address":            []interface{}{"localhost:6379"},
		"conn_max_idle_time": 300, // 300秒
	}
	AddMap("test_add_map_int", setting)

	opt, ok := options["test_add_map_int"]
	if !ok {
		t.Error("从map添加配置失败")
	}

	if opt.ConnMaxIdleTime != 300*time.Second {
		t.Errorf("期望ConnMaxIdleTime为300秒，实际为%v", opt.ConnMaxIdleTime)
	}
}

// TestAddMapBatch 测试批量添加配置
func TestAddMapBatch(t *testing.T) {
	batch := map[string]interface{}{
		"batch1": map[string]interface{}{
			"address": []interface{}{"localhost:6379"},
		},
		"batch2": map[string]interface{}{
			"address": []interface{}{"localhost:6380"},
		},
	}
	AddMapBatch(batch)

	if _, ok := options["batch1"]; !ok {
		t.Error("批量添加配置batch1失败")
	}
	if _, ok := options["batch2"]; !ok {
		t.Error("批量添加配置batch2失败")
	}
}

// TestUse 测试使用Redis连接
func TestUse(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	// 添加配置
	option := Option{
		Address: []string{s.Addr()},
	}
	Add("test_use", option)

	// 测试使用连接
	db := Use("test_use")
	if db == nil {
		t.Error("获取连接失败")
	}

	// 测试连接复用
	if db != Use("test_use") {
		t.Error("连接未能复用")
	}

	// 测试Ping
	if err := db.Ping(Ctx).Err(); err != nil {
		t.Errorf("Redis Ping失败: %v", err)
	}

	// 测试基本操作
	err := db.Set(Ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		t.Errorf("Redis Set失败: %v", err)
	}

	val, err := db.Get(Ctx, "test_key").Result()
	if err != nil {
		t.Errorf("Redis Get失败: %v", err)
	}
	if val != "test_value" {
		t.Errorf("期望获取test_value，实际获取%s", val)
	}

	// 测试不存在的键
	_, err = db.Get(Ctx, "the_key_does_not_exist_yeah").Result()
	if err != Nil {
		t.Errorf("Redis Get操作错误: %v", err)
	}
}

// TestUseDefault 测试使用默认连接
func TestUseDefault(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	// 添加默认配置
	option := Option{
		Address: []string{s.Addr()},
	}
	Add("default", option)

	// 空字符串应该使用default
	db := Use("")
	if db == nil {
		t.Error("获取默认连接失败")
	}

	if db != Use("default") {
		t.Error("空字符串未映射到default连接")
	}
}

// TestUseConcurrent 测试并发使用Redis连接
func TestUseConcurrent(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("test_concurrent", option)

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for loop := 0; loop < 10; loop++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			db := Use("test_concurrent")
			if db == nil {
				errors <- fmt.Errorf("goroutine %d: 获取连接失败", i)
				return
			}

			// 测试基本操作
			key := fmt.Sprintf("concurrent_key_%d", i)
			if err := db.Set(Ctx, key, i, 0).Err(); err != nil {
				errors <- fmt.Errorf("goroutine %d: Set失败: %v", i, err)
				return
			}

			val, err := db.Get(Ctx, key).Int()
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: Get失败: %v", i, err)
				return
			}

			if val != i {
				errors <- fmt.Errorf("goroutine %d: 期望值%d，实际值%d", i, i, val)
			}
		}(loop)
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	for err := range errors {
		t.Error(err)
	}
}

// TestUsePanic 测试使用不存在的配置
func TestUsePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望panic但没有发生")
		}
	}()

	Use("non_existent_config")
}

// TestAddPanic 测试添加空地址配置
func TestAddPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望panic但没有发生")
		}
	}()

	option := Option{
		Address: []string{}, // 空地址数组
	}
	Add("test_panic", option)
}
