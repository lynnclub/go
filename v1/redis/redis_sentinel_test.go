package redis

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestSentinelConfig 测试Redis Sentinel配置
func TestSentinelConfig(t *testing.T) {
	// 添加Sentinel配置
	option := Option{
		Address:      []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		Password:     "sentinel_password",
		MasterName:   "mymaster",
		PoolSize:     80,
		MinIdleConns: 8,
		MaxIdleConns: 30,
	}
	Add("test_sentinel_config", option)

	// 测试配置是否正确添加
	if opt, ok := options["test_sentinel_config"]; !ok {
		t.Error("Sentinel配置添加失败")
	} else {
		if len(opt.Address) != 3 {
			t.Errorf("期望3个哨兵地址，实际%d个", len(opt.Address))
		}
		if opt.Password != "sentinel_password" {
			t.Errorf("期望密码为sentinel_password，实际为%s", opt.Password)
		}
		if opt.MasterName != "mymaster" {
			t.Errorf("期望MasterName为mymaster，实际为%s", opt.MasterName)
		}
		if opt.PoolSize != 80 {
			t.Errorf("期望PoolSize为80，实际为%d", opt.PoolSize)
		}
	}
}

// TestSentinelWithCustomMasterName 测试自定义主库名称
func TestSentinelWithCustomMasterName(t *testing.T) {
	customMasterName := "custom_master"

	// 添加自定义主库名称的配置
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: customMasterName,
	}
	Add("test_sentinel_custom", option)

	if opt, ok := options["test_sentinel_custom"]; !ok {
		t.Error("自定义主库名称配置添加失败")
	} else if opt.MasterName != customMasterName {
		t.Errorf("期望主库名称为%s，实际为%s", customMasterName, opt.MasterName)
	}
}

// TestSentinelWithTLS 测试启用TLS的Sentinel配置
func TestSentinelWithTLS(t *testing.T) {
	// 添加启用TLS的Sentinel配置
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "mymaster",
		TLS:        true,
	}
	Add("test_sentinel_tls", option)

	if opt, ok := options["test_sentinel_tls"]; !ok {
		t.Error("TLS Sentinel配置添加失败")
	} else if !opt.TLS {
		t.Error("TLS配置未正确设置")
	}
}

// TestSentinelMultipleNodes 测试多哨兵节点配置
func TestSentinelMultipleNodes(t *testing.T) {
	sentinelAddrs := []string{
		"sentinel1:26379",
		"sentinel2:26379",
		"sentinel3:26379",
	}

	option := Option{
		Address:    sentinelAddrs,
		MasterName: "mymaster",
	}
	Add("test_sentinel_multi", option)

	if opt, ok := options["test_sentinel_multi"]; !ok {
		t.Error("多哨兵节点配置添加失败")
	} else {
		if len(opt.Address) != len(sentinelAddrs) {
			t.Errorf("期望%d个哨兵地址，实际%d个", len(sentinelAddrs), len(opt.Address))
		}
		for i, addr := range sentinelAddrs {
			if opt.Address[i] != addr {
				t.Errorf("哨兵地址%d期望%s，实际%s", i, addr, opt.Address[i])
			}
		}
	}
}

// TestSentinelDefaultValues 测试Sentinel配置默认值
func TestSentinelDefaultValues(t *testing.T) {
	option := Option{
		Address: []string{"localhost:26379"},
	}
	Add("test_sentinel_defaults", option)

	opt := options["test_sentinel_defaults"]

	if opt.PoolSize != 100 {
		t.Errorf("期望默认PoolSize为100，实际为%d", opt.PoolSize)
	}

	if opt.MasterName != "mymaster" {
		t.Errorf("期望默认MasterName为mymaster，实际为%s", opt.MasterName)
	}
}

// TestSentinelSingleNode 测试单哨兵节点配置（生产不推荐）
func TestSentinelSingleNode(t *testing.T) {
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "mymaster",
	}
	Add("test_sentinel_single", option)

	if opt, ok := options["test_sentinel_single"]; !ok {
		t.Error("单哨兵节点配置添加失败")
	} else if len(opt.Address) != 1 {
		t.Errorf("期望1个哨兵地址，实际%d个", len(opt.Address))
	}
}

// TestSentinelPoolSize 测试哨兵连接池大小配置
func TestSentinelPoolSize(t *testing.T) {
	tests := []struct {
		name     string
		poolSize int
		expected int
	}{
		{"自定义连接池大小30", 30, 30},
		{"自定义连接池大小150", 150, 150},
		{"自定义连接池大小1", 1, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_sentinel_pool_%d", tc.poolSize)
			option := Option{
				Address:    []string{"localhost:26379"},
				MasterName: "mymaster",
				PoolSize:   tc.poolSize,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.PoolSize != tc.expected {
				t.Errorf("期望PoolSize为%d，实际为%d", tc.expected, opt.PoolSize)
			}
		})
	}
}

// TestSentinelIdleConnections 测试哨兵空闲连接配置
func TestSentinelIdleConnections(t *testing.T) {
	option := Option{
		Address:      []string{"localhost:26379"},
		MasterName:   "mymaster",
		MinIdleConns: 15,
		MaxIdleConns: 60,
	}
	Add("test_sentinel_idle", option)

	opt := options["test_sentinel_idle"]
	if opt.MinIdleConns != 15 {
		t.Errorf("期望MinIdleConns为15，实际为%d", opt.MinIdleConns)
	}
	if opt.MaxIdleConns != 60 {
		t.Errorf("期望MaxIdleConns为60，实际为%d", opt.MaxIdleConns)
	}
}

// TestSentinelConnMaxIdleTime 测试哨兵连接最大空闲时间
func TestSentinelConnMaxIdleTime(t *testing.T) {
	tests := []struct {
		name     string
		idleTime time.Duration
		expected time.Duration
	}{
		{"默认空闲时间", 0, 5 * time.Minute},
		{"自定义15分钟", 15 * time.Minute, 15 * time.Minute},
		{"自定义1分钟", 1 * time.Minute, 1 * time.Minute},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_sentinel_idle_time_%v", tc.idleTime)
			option := Option{
				Address:         []string{"localhost:26379"},
				MasterName:      "mymaster",
				ConnMaxIdleTime: tc.idleTime,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.ConnMaxIdleTime != tc.expected {
				t.Errorf("期望ConnMaxIdleTime为%v，实际为%v", tc.expected, opt.ConnMaxIdleTime)
			}
		})
	}
}

// TestSentinelPasswordConfiguration 测试哨兵密码配置
func TestSentinelPasswordConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"无密码", ""},
		{"简单密码", "sentinel_pass"},
		{"复杂密码", "S3nt!n3l@P@ssw0rd#2024"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_sentinel_pwd_%s", tc.name)
			option := Option{
				Address:    []string{"localhost:26379"},
				MasterName: "mymaster",
				Password:   tc.password,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.Password != tc.password {
				t.Errorf("期望Password为'%s'，实际为'%s'", tc.password, opt.Password)
			}
		})
	}
}

// TestSentinelMasterNameVariations 测试不同的主库名称
func TestSentinelMasterNameVariations(t *testing.T) {
	tests := []struct {
		masterName string
	}{
		{"mymaster"},
		{"redis-master"},
		{"production-master"},
		{"master_01"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("MasterName_%s", tc.masterName), func(t *testing.T) {
			configName := fmt.Sprintf("test_sentinel_master_%s", tc.masterName)
			option := Option{
				Address:    []string{"localhost:26379"},
				MasterName: tc.masterName,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.MasterName != tc.masterName {
				t.Errorf("期望MasterName为%s，实际为%s", tc.masterName, opt.MasterName)
			}
		})
	}
}

// TestSentinelConfigFromMap 测试从map配置哨兵
func TestSentinelConfigFromMap(t *testing.T) {
	setting := map[string]interface{}{
		"address":            []interface{}{"localhost:26379", "localhost:26380", "localhost:26381"},
		"password":           "sentinel_map_password",
		"master_name":        "sentinel_master",
		"pool_size":          70,
		"min_idle_conns":     7,
		"max_idle_conns":     35,
		"conn_max_idle_time": "12m",
		"tls":                true,
	}
	AddMap("test_sentinel_from_map", setting)

	opt, ok := options["test_sentinel_from_map"]
	if !ok {
		t.Fatal("从map添加哨兵配置失败")
	}

	if len(opt.Address) != 3 {
		t.Errorf("期望3个哨兵地址，实际%d个", len(opt.Address))
	}
	if opt.Password != "sentinel_map_password" {
		t.Errorf("期望密码为sentinel_map_password，实际为%s", opt.Password)
	}
	if opt.MasterName != "sentinel_master" {
		t.Errorf("期望MasterName为sentinel_master，实际为%s", opt.MasterName)
	}
	if opt.PoolSize != 70 {
		t.Errorf("期望PoolSize为70，实际为%d", opt.PoolSize)
	}
	if opt.MinIdleConns != 7 {
		t.Errorf("期望MinIdleConns为7，实际为%d", opt.MinIdleConns)
	}
	if opt.MaxIdleConns != 35 {
		t.Errorf("期望MaxIdleConns为35，实际为%d", opt.MaxIdleConns)
	}
	if opt.ConnMaxIdleTime != 12*time.Minute {
		t.Errorf("期望ConnMaxIdleTime为12分钟，实际为%v", opt.ConnMaxIdleTime)
	}
	if !opt.TLS {
		t.Error("期望TLS为true，实际为false")
	}
}

// TestSentinelConcurrentConfigRead 测试并发读取哨兵配置
func TestSentinelConcurrentConfigRead(t *testing.T) {
	// 先添加配置
	for i := 0; i < 15; i++ {
		configName := fmt.Sprintf("concurrent_sentinel_%d", i)
		masterName := fmt.Sprintf("master_%d", i)
		option := Option{
			Address:    []string{fmt.Sprintf("localhost:%d", 26379+i)},
			MasterName: masterName,
			PoolSize:   60 + i,
		}
		Add(configName, option)
	}

	// 并发读取配置
	var wg sync.WaitGroup
	errors := make(chan error, 15)

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			configName := fmt.Sprintf("concurrent_sentinel_%d", idx)
			masterName := fmt.Sprintf("master_%d", idx)

			// 验证配置
			if opt, ok := options[configName]; !ok {
				errors <- fmt.Errorf("goroutine %d: 配置读取失败", idx)
			} else {
				if opt.MasterName != masterName {
					errors <- fmt.Errorf("goroutine %d: MasterName配置错误，期望%s，实际%s", idx, masterName, opt.MasterName)
				}
				if opt.PoolSize != 60+idx {
					errors <- fmt.Errorf("goroutine %d: PoolSize配置错误，期望%d，实际%d", idx, 60+idx, opt.PoolSize)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查错误
	for err := range errors {
		t.Error(err)
	}
}

// TestSentinelEmptyAddressPanic 测试空地址数组引发panic
func TestSentinelEmptyAddressPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望空地址数组引发panic但没有发生")
		}
	}()

	option := Option{
		Address:    []string{},
		MasterName: "mymaster",
	}
	Add("test_sentinel_empty_addr", option)
}

// TestSentinelConfigNotFoundPanic 测试使用不存在的配置引发panic
func TestSentinelConfigNotFoundPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望使用不存在的配置引发panic但没有发生")
		}
	}()

	Sentinel("non_existent_sentinel_config")
}

// TestSentinelHighAvailability 测试哨兵高可用配置（5个节点）
func TestSentinelHighAvailability(t *testing.T) {
	sentinelAddrs := []string{
		"sentinel1:26379",
		"sentinel2:26379",
		"sentinel3:26379",
		"sentinel4:26379",
		"sentinel5:26379",
	}

	option := Option{
		Address:    sentinelAddrs,
		MasterName: "ha_master",
	}
	Add("test_sentinel_ha", option)

	opt := options["test_sentinel_ha"]
	if len(opt.Address) != 5 {
		t.Errorf("期望5个哨兵地址，实际%d个", len(opt.Address))
	}
	if opt.MasterName != "ha_master" {
		t.Errorf("期望MasterName为ha_master，实际为%s", opt.MasterName)
	}
}

// TestSentinelWithDB 测试哨兵配置中的DB字段
func TestSentinelWithDB(t *testing.T) {
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "mymaster",
		DB:         3,
	}
	Add("test_sentinel_with_db", option)

	opt := options["test_sentinel_with_db"]
	if opt.DB != 3 {
		t.Errorf("期望DB为3，实际为%d", opt.DB)
	}
}

// TestSentinelEmptyMasterNameWithDefault 测试空主库名称使用默认值
func TestSentinelEmptyMasterNameWithDefault(t *testing.T) {
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "", // 空字符串，应该使用默认值
	}
	Add("test_sentinel_empty_master", option)

	opt := options["test_sentinel_empty_master"]
	if opt.MasterName != "mymaster" {
		t.Errorf("期望默认MasterName为mymaster，实际为%s", opt.MasterName)
	}
}
