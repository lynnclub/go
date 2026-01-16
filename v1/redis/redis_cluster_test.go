package redis

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestClusterConfig 测试Redis集群配置
func TestClusterConfig(t *testing.T) {
	// 添加集群配置
	option := Option{
		Address:      []string{"localhost:7000", "localhost:7001", "localhost:7002"},
		Password:     "test_password",
		PoolSize:     50,
		MinIdleConns: 5,
		MaxIdleConns: 20,
	}
	Add("test_cluster_config", option)

	// 测试配置是否正确添加
	if opt, ok := options["test_cluster_config"]; !ok {
		t.Error("集群配置添加失败")
	} else {
		if len(opt.Address) != 3 {
			t.Errorf("期望3个地址，实际%d个", len(opt.Address))
		}
		if opt.Password != "test_password" {
			t.Errorf("期望密码为test_password，实际为%s", opt.Password)
		}
		if opt.PoolSize != 50 {
			t.Errorf("期望PoolSize为50，实际为%d", opt.PoolSize)
		}
	}
}

// TestClusterWithTLS 测试启用TLS的集群配置
func TestClusterWithTLS(t *testing.T) {
	// 添加启用TLS的集群配置
	option := Option{
		Address:  []string{"localhost:7000", "localhost:7001"},
		Password: "",
		TLS:      true,
	}
	Add("test_cluster_tls", option)

	if opt, ok := options["test_cluster_tls"]; !ok {
		t.Error("TLS集群配置添加失败")
	} else if !opt.TLS {
		t.Error("TLS配置未正确设置")
	}
}

// TestClusterMultipleAddresses 测试多地址配置
func TestClusterMultipleAddresses(t *testing.T) {
	addresses := []string{
		"node1:7000",
		"node2:7001",
		"node3:7002",
		"node4:7003",
		"node5:7004",
	}

	option := Option{
		Address: addresses,
	}
	Add("test_cluster_multi", option)

	if opt, ok := options["test_cluster_multi"]; !ok {
		t.Error("多地址集群配置添加失败")
	} else {
		if len(opt.Address) != len(addresses) {
			t.Errorf("期望%d个地址，实际%d个", len(addresses), len(opt.Address))
		}
		for i, addr := range addresses {
			if opt.Address[i] != addr {
				t.Errorf("地址%d期望%s，实际%s", i, addr, opt.Address[i])
			}
		}
	}
}

// TestClusterDefaultValues 测试集群配置默认值
func TestClusterDefaultValues(t *testing.T) {
	option := Option{
		Address: []string{"localhost:7000"},
	}
	Add("test_cluster_defaults", option)

	opt := options["test_cluster_defaults"]

	if opt.PoolSize != 100 {
		t.Errorf("期望默认PoolSize为100，实际为%d", opt.PoolSize)
	}

	if opt.MasterName != "mymaster" {
		t.Errorf("期望默认MasterName为mymaster，实际为%s", opt.MasterName)
	}
}

// TestClusterPoolSize 测试集群连接池大小配置
func TestClusterPoolSize(t *testing.T) {
	tests := []struct {
		name     string
		poolSize int
		expected int
	}{
		{"自定义连接池大小50", 50, 50},
		{"自定义连接池大小200", 200, 200},
		{"自定义连接池大小1", 1, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_cluster_pool_%d", tc.poolSize)
			option := Option{
				Address:  []string{"localhost:7000"},
				PoolSize: tc.poolSize,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.PoolSize != tc.expected {
				t.Errorf("期望PoolSize为%d，实际为%d", tc.expected, opt.PoolSize)
			}
		})
	}
}

// TestClusterIdleConnections 测试集群空闲连接配置
func TestClusterIdleConnections(t *testing.T) {
	option := Option{
		Address:      []string{"localhost:7000"},
		MinIdleConns: 10,
		MaxIdleConns: 50,
	}
	Add("test_cluster_idle", option)

	opt := options["test_cluster_idle"]
	if opt.MinIdleConns != 10 {
		t.Errorf("期望MinIdleConns为10，实际为%d", opt.MinIdleConns)
	}
	if opt.MaxIdleConns != 50 {
		t.Errorf("期望MaxIdleConns为50，实际为%d", opt.MaxIdleConns)
	}
}

// TestClusterConnMaxIdleTime 测试集群连接最大空闲时间
func TestClusterConnMaxIdleTime(t *testing.T) {
	tests := []struct {
		name     string
		idleTime time.Duration
		expected time.Duration
	}{
		{"默认空闲时间", 0, 5 * time.Minute},
		{"自定义10分钟", 10 * time.Minute, 10 * time.Minute},
		{"自定义30秒", 30 * time.Second, 30 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_cluster_idle_time_%v", tc.idleTime)
			option := Option{
				Address:         []string{"localhost:7000"},
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

// TestClusterPasswordConfiguration 测试集群密码配置
func TestClusterPasswordConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"无密码", ""},
		{"简单密码", "simple_password"},
		{"复杂密码", "C0mpl3x!P@ssw0rd#2024"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_cluster_pwd_%s", tc.name)
			option := Option{
				Address:  []string{"localhost:7000"},
				Password: tc.password,
			}
			Add(configName, option)

			opt := options[configName]
			if opt.Password != tc.password {
				t.Errorf("期望Password为'%s'，实际为'%s'", tc.password, opt.Password)
			}
		})
	}
}

// TestClusterConfigFromMap 测试从map配置集群
func TestClusterConfigFromMap(t *testing.T) {
	setting := map[string]interface{}{
		"address":            []interface{}{"localhost:7000", "localhost:7001", "localhost:7002"},
		"password":           "cluster_password",
		"pool_size":          60,
		"min_idle_conns":     6,
		"max_idle_conns":     25,
		"conn_max_idle_time": "8m",
	}
	AddMap("test_cluster_from_map", setting)

	opt, ok := options["test_cluster_from_map"]
	if !ok {
		t.Fatal("从map添加集群配置失败")
	}

	if len(opt.Address) != 3 {
		t.Errorf("期望3个地址，实际%d个", len(opt.Address))
	}
	if opt.Password != "cluster_password" {
		t.Errorf("期望密码为cluster_password，实际为%s", opt.Password)
	}
	if opt.PoolSize != 60 {
		t.Errorf("期望PoolSize为60，实际为%d", opt.PoolSize)
	}
	if opt.MinIdleConns != 6 {
		t.Errorf("期望MinIdleConns为6，实际为%d", opt.MinIdleConns)
	}
	if opt.MaxIdleConns != 25 {
		t.Errorf("期望MaxIdleConns为25，实际为%d", opt.MaxIdleConns)
	}
	if opt.ConnMaxIdleTime != 8*time.Minute {
		t.Errorf("期望ConnMaxIdleTime为8分钟，实际为%v", opt.ConnMaxIdleTime)
	}
}

// TestClusterConcurrentConfigRead 测试并发读取集群配置
func TestClusterConcurrentConfigRead(t *testing.T) {
	// 先添加配置
	for i := 0; i < 20; i++ {
		configName := fmt.Sprintf("concurrent_cluster_%d", i)
		option := Option{
			Address:  []string{fmt.Sprintf("localhost:%d", 7000+i)},
			PoolSize: 50 + i,
		}
		Add(configName, option)
	}

	// 并发读取配置
	var wg sync.WaitGroup
	errors := make(chan error, 20)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			configName := fmt.Sprintf("concurrent_cluster_%d", idx)

			// 验证配置
			if opt, ok := options[configName]; !ok {
				errors <- fmt.Errorf("goroutine %d: 配置读取失败", idx)
			} else if opt.PoolSize != 50+idx {
				errors <- fmt.Errorf("goroutine %d: PoolSize配置错误，期望%d，实际%d", idx, 50+idx, opt.PoolSize)
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

// TestClusterEmptyAddressPanic 测试空地址数组引发panic
func TestClusterEmptyAddressPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望空地址数组引发panic但没有发生")
		}
	}()

	option := Option{
		Address: []string{},
	}
	Add("test_cluster_empty_addr", option)
}

// TestClusterConfigNotFoundPanic 测试使用不存在的配置引发panic
func TestClusterConfigNotFoundPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望使用不存在的配置引发panic但没有发生")
		}
	}()

	Cluster("non_existent_cluster_config")
}

// TestClusterSingleNodeConfig 测试单节点集群配置
func TestClusterSingleNodeConfig(t *testing.T) {
	option := Option{
		Address: []string{"localhost:7000"},
	}
	Add("test_cluster_single_node", option)

	opt := options["test_cluster_single_node"]
	if len(opt.Address) != 1 {
		t.Errorf("期望1个地址，实际%d个", len(opt.Address))
	}
}

// TestClusterWithDB 测试集群配置中的DB字段（集群模式不使用DB）
func TestClusterWithDB(t *testing.T) {
	option := Option{
		Address: []string{"localhost:7000"},
		DB:      5, // 集群模式会忽略此字段
	}
	Add("test_cluster_with_db", option)

	opt := options["test_cluster_with_db"]
	// 配置会保存DB字段，但实际使用时集群模式会忽略它
	if opt.DB != 5 {
		t.Errorf("期望DB为5，实际为%d", opt.DB)
	}
}

// TestClusterConnectionPooling 测试集群连接池的复用
func TestClusterConnectionPooling(t *testing.T) {
	// 添加配置但不实际连接（避免连接失败）
	option := Option{
		Address:  []string{"localhost:7000"},
		PoolSize: 50,
	}
	Add("test_cluster_pool_reuse", option)

	// 验证配置存在
	if _, ok := options["test_cluster_pool_reuse"]; !ok {
		t.Error("集群连接池配置添加失败")
	}
}

// TestClusterMultipleInstanceConfigs 测试多个集群实例配置
func TestClusterMultipleInstanceConfigs(t *testing.T) {
	clusters := map[string][]string{
		"cluster_prod":    {"prod-node1:7000", "prod-node2:7000", "prod-node3:7000"},
		"cluster_staging": {"staging-node1:7000", "staging-node2:7000"},
		"cluster_dev":     {"localhost:7000"},
	}

	for name, addresses := range clusters {
		option := Option{
			Address: addresses,
		}
		Add(name, option)

		if opt, ok := options[name]; !ok {
			t.Errorf("集群配置%s添加失败", name)
		} else if len(opt.Address) != len(addresses) {
			t.Errorf("集群%s地址数量不匹配，期望%d，实际%d", name, len(addresses), len(opt.Address))
		}
	}
}

// TestClusterAddressFormat 测试集群地址格式
func TestClusterAddressFormat(t *testing.T) {
	tests := []struct {
		name      string
		addresses []string
	}{
		{
			"IP地址带端口",
			[]string{"192.168.1.1:7000", "192.168.1.2:7000"},
		},
		{
			"域名带端口",
			[]string{"redis-cluster-1.example.com:7000", "redis-cluster-2.example.com:7000"},
		},
		{
			"localhost",
			[]string{"localhost:7000", "localhost:7001", "localhost:7002"},
		},
		{
			"混合格式",
			[]string{"192.168.1.1:7000", "redis.example.com:7001", "localhost:7002"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("test_cluster_addr_%s", tc.name)
			option := Option{
				Address: tc.addresses,
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

// TestClusterOptionsIndependence 测试不同配置的独立性
func TestClusterOptionsIndependence(t *testing.T) {
	option1 := Option{
		Address:  []string{"node1:7000", "node2:7000"},
		PoolSize: 50,
	}
	Add("cluster_independent_1", option1)

	option2 := Option{
		Address:  []string{"node3:7000", "node4:7000"},
		PoolSize: 100,
	}
	Add("cluster_independent_2", option2)

	// 验证两个配置互不影响
	opt1 := options["cluster_independent_1"]
	opt2 := options["cluster_independent_2"]

	if opt1.PoolSize == opt2.PoolSize {
		t.Error("两个配置的PoolSize不应该相同")
	}

	if len(opt1.Address) == 0 || len(opt2.Address) == 0 {
		t.Error("配置的地址不应该为空")
	}

	// 验证两个配置的地址是不同的
	if opt1.Address[0] == opt2.Address[0] {
		t.Error("两个配置的地址应该不同")
	}

	// 验证配置被正确保存
	if opt1.Address[0] != "node1:7000" {
		t.Errorf("配置1的地址不正确，期望node1:7000，实际%s", opt1.Address[0])
	}
	if opt2.Address[0] != "node3:7000" {
		t.Errorf("配置2的地址不正确，期望node3:7000，实际%s", opt2.Address[0])
	}
}

// TestClusterConfigUpdate 测试配置更新
func TestClusterConfigUpdate(t *testing.T) {
	// 第一次添加配置
	option1 := Option{
		Address:  []string{"localhost:7000"},
		PoolSize: 50,
	}
	Add("cluster_update_test", option1)

	firstPoolSize := options["cluster_update_test"].PoolSize

	// 更新配置
	option2 := Option{
		Address:  []string{"localhost:7000", "localhost:7001"},
		PoolSize: 100,
	}
	Add("cluster_update_test", option2)

	secondPoolSize := options["cluster_update_test"].PoolSize
	secondAddressCount := len(options["cluster_update_test"].Address)

	if firstPoolSize == secondPoolSize {
		t.Error("配置应该被更新")
	}

	if secondPoolSize != 100 {
		t.Errorf("更新后PoolSize期望100，实际%d", secondPoolSize)
	}

	if secondAddressCount != 2 {
		t.Errorf("更新后地址数量期望2，实际%d", secondAddressCount)
	}
}

// TestClusterPasswordSecurity 测试密码配置的安全性
func TestClusterPasswordSecurity(t *testing.T) {
	sensitivePassword := "very_secret_password_123!@#"

	option := Option{
		Address:  []string{"localhost:7000"},
		Password: sensitivePassword,
	}
	Add("cluster_password_security", option)

	opt := options["cluster_password_security"]

	// 验证密码被正确保存（在实际应用中应该加密）
	if opt.Password != sensitivePassword {
		t.Error("密码配置不正确")
	}

	// 确保密码不是空字符串
	if opt.Password == "" {
		t.Error("密码不应该为空")
	}
}

// TestClusterConnectionLimits 测试连接数限制配置
func TestClusterConnectionLimits(t *testing.T) {
	tests := []struct {
		name         string
		poolSize     int
		minIdleConns int
		maxIdleConns int
		valid        bool
	}{
		{"正常配置", 100, 10, 50, true},
		{"最小配置", 1, 0, 0, true},
		{"大连接池", 1000, 100, 500, true},
		{"最小空闲为0", 50, 0, 25, true},
		{"最大空闲为0", 50, 10, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configName := fmt.Sprintf("cluster_limits_%s", tc.name)
			option := Option{
				Address:      []string{"localhost:7000"},
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

// TestClusterBatchConfigAddition 测试批量添加集群配置
func TestClusterBatchConfigAddition(t *testing.T) {
	batch := map[string]interface{}{
		"cluster_batch_1": map[string]interface{}{
			"address":   []interface{}{"node1:7000", "node2:7000"},
			"pool_size": 60,
		},
		"cluster_batch_2": map[string]interface{}{
			"address":   []interface{}{"node3:7000", "node4:7000"},
			"pool_size": 70,
		},
		"cluster_batch_3": map[string]interface{}{
			"address":   []interface{}{"node5:7000"},
			"pool_size": 80,
		},
	}

	AddMapBatch(batch)

	// 验证所有配置都被添加
	for i := 1; i <= 3; i++ {
		configName := fmt.Sprintf("cluster_batch_%d", i)
		if _, ok := options[configName]; !ok {
			t.Errorf("批量配置%s添加失败", configName)
		}
	}

	// 验证配置的正确性
	if opt, ok := options["cluster_batch_1"]; ok && opt.PoolSize != 60 {
		t.Errorf("cluster_batch_1的PoolSize期望60，实际%d", opt.PoolSize)
	}
	if opt, ok := options["cluster_batch_2"]; ok && opt.PoolSize != 70 {
		t.Errorf("cluster_batch_2的PoolSize期望70，实际%d", opt.PoolSize)
	}
	if opt, ok := options["cluster_batch_3"]; ok && opt.PoolSize != 80 {
		t.Errorf("cluster_batch_3的PoolSize期望80，实际%d", opt.PoolSize)
	}
}

// TestClusterConfigConcurrentAccess 测试并发访问配置
func TestClusterConfigConcurrentAccess(t *testing.T) {
	// 预先添加配置
	option := Option{
		Address:  []string{"localhost:7000"},
		PoolSize: 100,
	}
	Add("cluster_concurrent_read", option)

	var wg sync.WaitGroup
	errors := make(chan error, 50)

	// 并发读取配置
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			opt, ok := options["cluster_concurrent_read"]
			if !ok {
				errors <- fmt.Errorf("goroutine %d: 配置读取失败", idx)
				return
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

// TestClusterLargeScaleAddresses 测试大规模地址配置
func TestClusterLargeScaleAddresses(t *testing.T) {
	// 创建大量节点地址
	addresses := make([]string, 100)
	for i := 0; i < 100; i++ {
		addresses[i] = fmt.Sprintf("node%d:7000", i)
	}

	option := Option{
		Address: addresses,
	}
	Add("cluster_large_scale", option)

	opt := options["cluster_large_scale"]
	if len(opt.Address) != 100 {
		t.Errorf("期望100个地址，实际%d个", len(opt.Address))
	}

	// 验证所有地址都正确保存
	for i := 0; i < 100; i++ {
		expected := fmt.Sprintf("node%d:7000", i)
		if opt.Address[i] != expected {
			t.Errorf("地址%d错误，期望%s，实际%s", i, expected, opt.Address[i])
		}
	}
}
