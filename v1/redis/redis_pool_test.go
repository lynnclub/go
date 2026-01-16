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

// TestUseWithDB 测试使用不同的DB
func TestUseWithDB(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	// 为不同的DB配置连接
	for db := 0; db < 3; db++ {
		configName := fmt.Sprintf("test_db_%d", db)
		option := Option{
			Address: []string{s.Addr()},
			DB:      db,
		}
		Add(configName, option)
	}

	// 测试不同DB的数据隔离
	client0 := Use("test_db_0")
	client1 := Use("test_db_1")

	// 在DB 0中设置值
	err := client0.Set(Ctx, "isolation_test", "db0_value", 0).Err()
	if err != nil {
		t.Errorf("DB 0 Set失败: %v", err)
	}

	// 在DB 1中设置相同键但不同值
	err = client1.Set(Ctx, "isolation_test", "db1_value", 0).Err()
	if err != nil {
		t.Errorf("DB 1 Set失败: %v", err)
	}

	// 验证数据隔离
	val0, _ := client0.Get(Ctx, "isolation_test").Result()
	val1, _ := client1.Get(Ctx, "isolation_test").Result()

	if val0 != "db0_value" {
		t.Errorf("DB 0期望值db0_value，实际%s", val0)
	}
	if val1 != "db1_value" {
		t.Errorf("DB 1期望值db1_value，实际%s", val1)
	}
}

// TestUseWithTLS 测试TLS配置
func TestUseWithTLS(t *testing.T) {
	option := Option{
		Address: []string{"localhost:6380"},
		TLS:     true,
	}
	Add("test_tls_config", option)

	// 验证TLS配置已保存
	opt := options["test_tls_config"]
	if !opt.TLS {
		t.Error("TLS配置未正确设置")
	}
}

// TestAddWithAllOptions 测试设置所有选项
func TestAddWithAllOptions(t *testing.T) {
	option := Option{
		Address:         []string{"localhost:6379"},
		Password:        "complex_password",
		DB:              5,
		PoolSize:        200,
		MinIdleConns:    20,
		MaxIdleConns:    100,
		ConnMaxIdleTime: 10 * time.Minute,
		MasterName:      "custom_master",
		TLS:             true,
	}
	Add("test_all_options", option)

	opt := options["test_all_options"]
	if opt.Password != "complex_password" {
		t.Errorf("Password配置错误，期望complex_password，实际%s", opt.Password)
	}
	if opt.DB != 5 {
		t.Errorf("DB配置错误，期望5，实际%d", opt.DB)
	}
	if opt.PoolSize != 200 {
		t.Errorf("PoolSize配置错误，期望200，实际%d", opt.PoolSize)
	}
	if opt.MinIdleConns != 20 {
		t.Errorf("MinIdleConns配置错误，期望20，实际%d", opt.MinIdleConns)
	}
	if opt.MaxIdleConns != 100 {
		t.Errorf("MaxIdleConns配置错误，期望100，实际%d", opt.MaxIdleConns)
	}
	if opt.ConnMaxIdleTime != 10*time.Minute {
		t.Errorf("ConnMaxIdleTime配置错误，期望10m，实际%v", opt.ConnMaxIdleTime)
	}
	if opt.MasterName != "custom_master" {
		t.Errorf("MasterName配置错误，期望custom_master，实际%s", opt.MasterName)
	}
	if !opt.TLS {
		t.Error("TLS配置错误，期望true，实际false")
	}
}

// TestRedisOperations 测试各种Redis操作
func TestRedisOperations(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("test_operations", option)

	client := Use("test_operations")

	// 测试String操作
	t.Run("String操作", func(t *testing.T) {
		err := client.Set(Ctx, "string_key", "string_value", 0).Err()
		if err != nil {
			t.Errorf("Set失败: %v", err)
		}

		val, err := client.Get(Ctx, "string_key").Result()
		if err != nil {
			t.Errorf("Get失败: %v", err)
		}
		if val != "string_value" {
			t.Errorf("期望string_value，实际%s", val)
		}

		// 测试删除
		err = client.Del(Ctx, "string_key").Err()
		if err != nil {
			t.Errorf("Del失败: %v", err)
		}

		_, err = client.Get(Ctx, "string_key").Result()
		if err != Nil {
			t.Errorf("删除后应该返回Nil错误，实际: %v", err)
		}
	})

	// 测试Hash操作
	t.Run("Hash操作", func(t *testing.T) {
		err := client.HSet(Ctx, "hash_key", "field1", "value1").Err()
		if err != nil {
			t.Errorf("HSet失败: %v", err)
		}

		val, err := client.HGet(Ctx, "hash_key", "field1").Result()
		if err != nil {
			t.Errorf("HGet失败: %v", err)
		}
		if val != "value1" {
			t.Errorf("期望value1，实际%s", val)
		}
	})

	// 测试List操作
	t.Run("List操作", func(t *testing.T) {
		err := client.RPush(Ctx, "list_key", "item1", "item2", "item3").Err()
		if err != nil {
			t.Errorf("RPush失败: %v", err)
		}

		length, err := client.LLen(Ctx, "list_key").Result()
		if err != nil {
			t.Errorf("LLen失败: %v", err)
		}
		if length != 3 {
			t.Errorf("期望列表长度3，实际%d", length)
		}

		val, err := client.LIndex(Ctx, "list_key", 0).Result()
		if err != nil {
			t.Errorf("LIndex失败: %v", err)
		}
		if val != "item1" {
			t.Errorf("期望item1，实际%s", val)
		}
	})

	// 测试Set操作
	t.Run("Set操作", func(t *testing.T) {
		err := client.SAdd(Ctx, "set_key", "member1", "member2", "member3").Err()
		if err != nil {
			t.Errorf("SAdd失败: %v", err)
		}

		count, err := client.SCard(Ctx, "set_key").Result()
		if err != nil {
			t.Errorf("SCard失败: %v", err)
		}
		if count != 3 {
			t.Errorf("期望集合大小3，实际%d", count)
		}

		isMember, err := client.SIsMember(Ctx, "set_key", "member1").Result()
		if err != nil {
			t.Errorf("SIsMember失败: %v", err)
		}
		if !isMember {
			t.Error("member1应该在集合中")
		}
	})

	// 测试过期时间
	t.Run("过期时间", func(t *testing.T) {
		err := client.Set(Ctx, "expire_key", "value", 1*time.Second).Err()
		if err != nil {
			t.Errorf("Set with expire失败: %v", err)
		}

		ttl, err := client.TTL(Ctx, "expire_key").Result()
		if err != nil {
			t.Errorf("TTL失败: %v", err)
		}
		if ttl <= 0 {
			t.Errorf("TTL应该大于0，实际%v", ttl)
		}
	})
}

// TestAddMapWithMissingFields 测试从map添加配置（缺少可选字段）
func TestAddMapWithMissingFields(t *testing.T) {
	setting := map[string]interface{}{
		"address": []interface{}{"localhost:6379"},
		// 所有其他字段都缺失，应该使用默认值
	}
	AddMap("test_map_minimal", setting)

	opt, ok := options["test_map_minimal"]
	if !ok {
		t.Fatal("配置添加失败")
	}

	// 验证默认值
	if opt.PoolSize != 100 {
		t.Errorf("期望默认PoolSize为100，实际%d", opt.PoolSize)
	}
	if opt.MasterName != "mymaster" {
		t.Errorf("期望默认MasterName为mymaster，实际%s", opt.MasterName)
	}
	if opt.ConnMaxIdleTime != 5*time.Minute {
		t.Errorf("期望默认ConnMaxIdleTime为5分钟，实际%v", opt.ConnMaxIdleTime)
	}
}

// TestNilError 测试Nil错误常量
func TestNilError(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("test_nil_error", option)

	client := Use("test_nil_error")

	// 获取不存在的键应该返回Nil错误
	_, err := client.Get(Ctx, "non_existent_key_for_nil_test").Result()
	if err != Nil {
		t.Errorf("期望Nil错误，实际%v", err)
	}
}

// TestContextUsage 测试Context的使用
func TestContextUsage(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	option := Option{
		Address: []string{s.Addr()},
	}
	Add("test_context", option)

	client := Use("test_context")

	// 使用全局Ctx
	err := client.Set(Ctx, "ctx_key", "ctx_value", 0).Err()
	if err != nil {
		t.Errorf("使用Ctx Set失败: %v", err)
	}

	val, err := client.Get(Ctx, "ctx_key").Result()
	if err != nil {
		t.Errorf("使用Ctx Get失败: %v", err)
	}
	if val != "ctx_value" {
		t.Errorf("期望ctx_value，实际%s", val)
	}
}

// TestMultipleConfigs 测试多个不同配置共存
func TestMultipleConfigs(t *testing.T) {
	s1 := setupMiniRedis(t)
	defer s1.Close()

	s2 := setupMiniRedis(t)
	defer s2.Close()

	// 添加两个不同的配置
	Add("redis1", Option{Address: []string{s1.Addr()}, DB: 0})
	Add("redis2", Option{Address: []string{s2.Addr()}, DB: 1})

	client1 := Use("redis1")
	client2 := Use("redis2")

	// 在两个不同的Redis实例中设置相同的键
	client1.Set(Ctx, "multi_config_key", "value1", 0)
	client2.Set(Ctx, "multi_config_key", "value2", 0)

	val1, _ := client1.Get(Ctx, "multi_config_key").Result()
	val2, _ := client2.Get(Ctx, "multi_config_key").Result()

	if val1 != "value1" {
		t.Errorf("redis1期望value1，实际%s", val1)
	}
	if val2 != "value2" {
		t.Errorf("redis2期望value2，实际%s", val2)
	}
}

// TestPoolSizeConfiguration 测试连接池大小配置的影响
func TestPoolSizeConfiguration(t *testing.T) {
	s := setupMiniRedis(t)
	defer s.Close()

	tests := []struct {
		name     string
		poolSize int
	}{
		{"小连接池", 1},
		{"中连接池", 50},
		{"大连接池", 200},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("pool_size_test_%d", tc.poolSize)
			option := Option{
				Address:  []string{s.Addr()},
				PoolSize: tc.poolSize,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.PoolSize != tc.poolSize {
				t.Errorf("期望PoolSize为%d，实际%d", tc.poolSize, opt.PoolSize)
			}

			// 测试连接可用性
			client := Use(configName)
			err := client.Ping(Ctx).Err()
			if err != nil {
				t.Errorf("Ping失败: %v", err)
			}
		})
	}
}
