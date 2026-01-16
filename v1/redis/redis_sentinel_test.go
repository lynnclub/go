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

// TestSentinelConnectionPooling 测试哨兵连接池的复用
func TestSentinelConnectionPooling(t *testing.T) {
	// 添加配置但不实际连接（避免连接失败）
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "mymaster",
		PoolSize:   50,
	}
	Add("test_sentinel_pool_reuse", option)

	// 验证配置存在
	if _, ok := options["test_sentinel_pool_reuse"]; !ok {
		t.Error("哨兵连接池配置添加失败")
	}
}

// TestSentinelMultipleInstanceConfigs 测试多个哨兵实例配置
func TestSentinelMultipleInstanceConfigs(t *testing.T) {
	sentinels := map[string]struct {
		addresses  []string
		masterName string
	}{
		"sentinel_prod": {
			[]string{"prod-sentinel1:26379", "prod-sentinel2:26379", "prod-sentinel3:26379"},
			"prod-master",
		},
		"sentinel_staging": {
			[]string{"staging-sentinel1:26379", "staging-sentinel2:26379"},
			"staging-master",
		},
		"sentinel_dev": {
			[]string{"localhost:26379"},
			"dev-master",
		},
	}

	for name, config := range sentinels {
		option := Option{
			Address:    config.addresses,
			MasterName: config.masterName,
		}
		Add(name, option)

		if opt, ok := options[name]; !ok {
			t.Errorf("哨兵配置%s添加失败", name)
		} else {
			if len(opt.Address) != len(config.addresses) {
				t.Errorf("哨兵%s地址数量不匹配，期望%d，实际%d", name, len(config.addresses), len(opt.Address))
			}
			if opt.MasterName != config.masterName {
				t.Errorf("哨兵%s主库名称不匹配，期望%s，实际%s", name, config.masterName, opt.MasterName)
			}
		}
	}
}

// TestSentinelAddressFormat 测试哨兵地址格式
func TestSentinelAddressFormat(t *testing.T) {
	tests := []struct {
		name      string
		addresses []string
	}{
		{
			"IP地址带端口",
			[]string{"192.168.1.10:26379", "192.168.1.11:26379"},
		},
		{
			"域名带端口",
			[]string{"sentinel-1.example.com:26379", "sentinel-2.example.com:26379"},
		},
		{
			"localhost",
			[]string{"localhost:26379", "localhost:26380", "localhost:26381"},
		},
		{
			"混合格式",
			[]string{"192.168.1.10:26379", "sentinel.example.com:26379", "localhost:26379"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_sentinel_addr_%s", tc.name)
			option := Option{
				Address:    tc.addresses,
				MasterName: "mymaster",
			}
			Add(configName, option)

			opt := options[configName]
			for i, addr := range tc.addresses {
				if opt.Address[i] != addr {
					t.Errorf("地址%d不匹配，期望%s，实际%s", i, addr, opt.Address[i])
				}
			}
		})
	}
}

// TestSentinelOptionsIndependence 测试不同配置的独立性
func TestSentinelOptionsIndependence(t *testing.T) {
	option1 := Option{
		Address:    []string{"sentinel1:26379", "sentinel2:26379"},
		MasterName: "master1",
		PoolSize:   50,
	}
	Add("sentinel_independent_1", option1)

	option2 := Option{
		Address:    []string{"sentinel3:26379", "sentinel4:26379"},
		MasterName: "master2",
		PoolSize:   100,
	}
	Add("sentinel_independent_2", option2)

	// 验证两个配置互不影响
	opt1 := options["sentinel_independent_1"]
	opt2 := options["sentinel_independent_2"]

	if opt1.MasterName == opt2.MasterName {
		t.Error("两个配置的MasterName不应该相同")
	}

	if opt1.PoolSize == opt2.PoolSize {
		t.Error("两个配置的PoolSize不应该相同")
	}

	// 验证两个配置的地址是不同的
	if opt1.Address[0] == opt2.Address[0] {
		t.Error("两个配置的地址应该不同")
	}

	// 验证配置被正确保存
	if opt1.Address[0] != "sentinel1:26379" {
		t.Errorf("配置1的地址不正确，期望sentinel1:26379，实际%s", opt1.Address[0])
	}
	if opt2.Address[0] != "sentinel3:26379" {
		t.Errorf("配置2的地址不正确，期望sentinel3:26379，实际%s", opt2.Address[0])
	}
}

// TestSentinelConfigUpdate 测试配置更新
func TestSentinelConfigUpdate(t *testing.T) {
	// 第一次添加配置
	option1 := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "old_master",
		PoolSize:   50,
	}
	Add("sentinel_update_test", option1)

	firstMasterName := options["sentinel_update_test"].MasterName

	// 更新配置
	option2 := Option{
		Address:    []string{"localhost:26379", "localhost:26380"},
		MasterName: "new_master",
		PoolSize:   100,
	}
	Add("sentinel_update_test", option2)

	secondMasterName := options["sentinel_update_test"].MasterName
	secondPoolSize := options["sentinel_update_test"].PoolSize
	secondAddressCount := len(options["sentinel_update_test"].Address)

	if firstMasterName == secondMasterName {
		t.Error("MasterName应该被更新")
	}

	if secondMasterName != "new_master" {
		t.Errorf("更新后MasterName期望new_master，实际%s", secondMasterName)
	}

	if secondPoolSize != 100 {
		t.Errorf("更新后PoolSize期望100，实际%d", secondPoolSize)
	}

	if secondAddressCount != 2 {
		t.Errorf("更新后地址数量期望2，实际%d", secondAddressCount)
	}
}

// TestSentinelPasswordSecurity 测试密码配置的安全性
func TestSentinelPasswordSecurity(t *testing.T) {
	sensitivePassword := "sentinel_secret_password_456!@#"

	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "mymaster",
		Password:   sensitivePassword,
	}
	Add("sentinel_password_security", option)

	opt := options["sentinel_password_security"]

	// 验证密码被正确保存
	if opt.Password != sensitivePassword {
		t.Error("密码配置不正确")
	}
}

// TestSentinelConnectionLimits 测试连接数限制配置
func TestSentinelConnectionLimits(t *testing.T) {
	tests := []struct {
		name         string
		poolSize     int
		minIdleConns int
		maxIdleConns int
		valid        bool
	}{
		{"正常配置", 80, 8, 40, true},
		{"最小配置", 1, 0, 0, true},
		{"大连接池", 500, 50, 250, true},
		{"最小空闲为0", 60, 0, 30, true},
		{"最大空闲为0", 60, 10, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("sentinel_limits_%s", tc.name)
			option := Option{
				Address:      []string{"localhost:26379"},
				MasterName:   "mymaster",
				PoolSize:     tc.poolSize,
				MinIdleConns: tc.minIdleConns,
				MaxIdleConns: tc.maxIdleConns,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.PoolSize != tc.poolSize {
				t.Errorf("PoolSize期望%d，实际%d", tc.poolSize, opt.PoolSize)
			}
			if opt.MinIdleConns != tc.minIdleConns {
				t.Errorf("MinIdleConns期望%d，实际%d", tc.minIdleConns, opt.MinIdleConns)
			}
			if opt.MaxIdleConns != tc.maxIdleConns {
				t.Errorf("MaxIdleConns期望%d，实际%d", tc.maxIdleConns, opt.MaxIdleConns)
			}
		})
	}
}

// TestSentinelBatchConfigAddition 测试批量添加哨兵配置
func TestSentinelBatchConfigAddition(t *testing.T) {
	batch := map[string]interface{}{
		"sentinel_batch_1": map[string]interface{}{
			"address":     []interface{}{"sentinel1:26379", "sentinel2:26379"},
			"master_name": "batch_master_1",
			"pool_size":   60,
		},
		"sentinel_batch_2": map[string]interface{}{
			"address":     []interface{}{"sentinel3:26379", "sentinel4:26379"},
			"master_name": "batch_master_2",
			"pool_size":   70,
		},
		"sentinel_batch_3": map[string]interface{}{
			"address":     []interface{}{"sentinel5:26379"},
			"master_name": "batch_master_3",
			"pool_size":   80,
		},
	}

	AddMapBatch(batch)

	// 验证所有配置都被添加
	for i := 1; i <= 3; i++ {
		configName := fmt.Sprintf("sentinel_batch_%d", i)
		if _, ok := options[configName]; !ok {
			t.Errorf("批量配置%s添加失败", configName)
		}
	}

	// 验证配置的正确性
	if opt, ok := options["sentinel_batch_1"]; ok {
		if opt.PoolSize != 60 {
			t.Errorf("sentinel_batch_1的PoolSize期望60，实际%d", opt.PoolSize)
		}
		if opt.MasterName != "batch_master_1" {
			t.Errorf("sentinel_batch_1的MasterName期望batch_master_1，实际%s", opt.MasterName)
		}
	}
}

// TestSentinelConfigConcurrentAccess 测试并发访问配置
func TestSentinelConfigConcurrentAccess(t *testing.T) {
	// 预先添加配置
	option := Option{
		Address:    []string{"localhost:26379"},
		MasterName: "concurrent_master",
		PoolSize:   100,
	}
	Add("sentinel_concurrent_read", option)

	var wg sync.WaitGroup
	errors := make(chan error, 50)

	// 并发读取配置
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			opt, ok := options["sentinel_concurrent_read"]
			if !ok {
				errors <- fmt.Errorf("goroutine %d: 配置读取失败", idx)
				return
			}

			if opt.MasterName != "concurrent_master" {
				errors <- fmt.Errorf("goroutine %d: MasterName错误，期望concurrent_master，实际%s", idx, opt.MasterName)
			}

			if opt.PoolSize != 100 {
				errors <- fmt.Errorf("goroutine %d: PoolSize错误，期望100，实际%d", idx, opt.PoolSize)
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

// TestSentinelLargeScaleAddresses 测试大规模哨兵地址配置
func TestSentinelLargeScaleAddresses(t *testing.T) {
	// 创建大量哨兵节点地址（虽然实际不会有这么多）
	addresses := make([]string, 20)
	for i := 0; i < 20; i++ {
		addresses[i] = fmt.Sprintf("sentinel-node%d:26379", i)
	}

	option := Option{
		Address:    addresses,
		MasterName: "large_scale_master",
	}
	Add("sentinel_large_scale", option)

	opt := options["sentinel_large_scale"]
	if len(opt.Address) != 20 {
		t.Errorf("期望20个哨兵地址，实际%d个", len(opt.Address))
	}

	// 验证所有地址都正确保存
	for i := 0; i < 20; i++ {
		expected := fmt.Sprintf("sentinel-node%d:26379", i)
		if opt.Address[i] != expected {
			t.Errorf("地址%d错误，期望%s，实际%s", i, expected, opt.Address[i])
		}
	}
}

// TestSentinelFailoverOptions 测试故障转移相关配置
func TestSentinelFailoverOptions(t *testing.T) {
	tests := []struct {
		name            string
		masterName      string
		addresses       []string
		connMaxIdleTime time.Duration
	}{
		{
			"短空闲时间",
			"master1",
			[]string{"sentinel1:26379", "sentinel2:26379"},
			30 * time.Second,
		},
		{
			"中等空闲时间",
			"master2",
			[]string{"sentinel3:26379", "sentinel4:26379"},
			5 * time.Minute,
		},
		{
			"长空闲时间",
			"master3",
			[]string{"sentinel5:26379", "sentinel6:26379"},
			30 * time.Minute,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("sentinel_failover_%s", tc.name)
			option := Option{
				Address:         tc.addresses,
				MasterName:      tc.masterName,
				ConnMaxIdleTime: tc.connMaxIdleTime,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.MasterName != tc.masterName {
				t.Errorf("MasterName期望%s，实际%s", tc.masterName, opt.MasterName)
			}
			if opt.ConnMaxIdleTime != tc.connMaxIdleTime {
				t.Errorf("ConnMaxIdleTime期望%v，实际%v", tc.connMaxIdleTime, opt.ConnMaxIdleTime)
			}
		})
	}
}

// TestSentinelOddEvenNodeCounts 测试奇数和偶数个哨兵节点
func TestSentinelOddEvenNodeCounts(t *testing.T) {
	tests := []struct {
		name      string
		nodeCount int
	}{
		{"单节点", 1},
		{"双节点", 2},
		{"三节点(推荐)", 3},
		{"五节点", 5},
		{"七节点", 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addresses := make([]string, tc.nodeCount)
			for i := 0; i < tc.nodeCount; i++ {
				addresses[i] = fmt.Sprintf("sentinel%d:26379", i)
			}

			configName := fmt.Sprintf("sentinel_nodes_%d", tc.nodeCount)
			option := Option{
				Address:    addresses,
				MasterName: "mymaster",
			}
			Add(configName, option)

			opt := options[configName]
			if len(opt.Address) != tc.nodeCount {
				t.Errorf("期望%d个节点，实际%d个", tc.nodeCount, len(opt.Address))
			}
		})
	}
}

// TestSentinelWithDBConfiguration 测试哨兵模式中的DB配置
func TestSentinelWithDBConfiguration(t *testing.T) {
	// 哨兵模式可以指定DB
	tests := []struct {
		db int
	}{
		{0},
		{1},
		{5},
		{15},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("DB_%d", tc.db), func(t *testing.T) {
			configName := fmt.Sprintf("sentinel_db_%d", tc.db)
			option := Option{
				Address:    []string{"localhost:26379"},
				MasterName: "mymaster",
				DB:         tc.db,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.DB != tc.db {
				t.Errorf("DB期望%d，实际%d", tc.db, opt.DB)
			}
		})
	}
}

// TestSentinelMasterNameSpecialCharacters 测试主库名称中的特殊字符
func TestSentinelMasterNameSpecialCharacters(t *testing.T) {
	tests := []struct {
		masterName string
	}{
		{"master-with-dash"},
		{"master_with_underscore"},
		{"master.with.dot"},
		{"master123"},
		{"MASTER_UPPERCASE"},
		{"MixedCase_Master-123"},
	}

	for _, tc := range tests {
		t.Run(tc.masterName, func(t *testing.T) {
			configName := fmt.Sprintf("sentinel_special_%s", tc.masterName)
			option := Option{
				Address:    []string{"localhost:26379"},
				MasterName: tc.masterName,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.MasterName != tc.masterName {
				t.Errorf("MasterName期望%s，实际%s", tc.masterName, opt.MasterName)
			}
		})
	}
}
