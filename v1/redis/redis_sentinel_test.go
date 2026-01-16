package redis

import (
	"testing"
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

// 注意：实际的Sentinel连接测试需要真实的Redis Sentinel环境
// 在单元测试中，我们主要测试配置管理部分
// 集成测试应在有真实Sentinel环境时进行
