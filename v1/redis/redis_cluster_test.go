package redis

import (
	"testing"
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

// 注意：实际的Cluster连接测试需要真实的Redis集群环境
// 在单元测试中，我们主要测试配置管理部分
// 集成测试应在有真实集群环境时进行
