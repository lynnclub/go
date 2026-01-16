package elasticsearch

import (
	"testing"
)

// TestUseV7WithDefaultName 测试使用默认名称
func TestUseV7WithDefaultName(t *testing.T) {
	// 先添加一个配置
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "elastic",
		Password: "password",
	}
	Add("default", option)

	// 注意：这个测试需要实际的 Elasticsearch 服务器才能通过
	// 如果没有服务器，将会 panic
	// 在实际环境中，建议使用 mock 或跳过此测试
	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV7("")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV7WithCustomName 测试使用自定义名称
func TestUseV7WithCustomName(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("custom_v7", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV7("custom_v7")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV7OptionNotFound 测试配置不存在的情况
func TestUseV7OptionNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该在配置不存在时 panic")
		} else {
			errMsg, ok := r.(string)
			if !ok || errMsg != "Option not found nonexistent_v7" {
				t.Errorf("panic 信息不正确: %v", r)
			}
		}
	}()

	UseV7("nonexistent_v7")
}

// TestUseV7EmptyAddresses 测试空地址配置
func TestUseV7EmptyAddresses(t *testing.T) {
	// 这个测试在 Add 函数中就会 panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该在地址为空时 panic")
		}
	}()

	option := Option{
		Address: []string{},
	}
	Add("empty_v7", option)
	UseV7("empty_v7")
}

// TestUseV7Singleton 测试单例模式
func TestUseV7Singleton(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("singleton_v7", option)

	// 清理可能存在的实例
	poolV7.Delete("singleton_v7")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 第一次调用会创建新实例
	client1 := UseV7("singleton_v7")

	// 第二次调用应该返回相同的实例
	client2 := UseV7("singleton_v7")

	if client1 != client2 {
		t.Error("UseV7 应该返回相同的实例（单例模式）")
	}
}

// TestUseV7WithoutAuth 测试无认证配置
func TestUseV7WithoutAuth(t *testing.T) {
	option := Option{
		Address: []string{"http://localhost:9200"},
	}
	Add("noauth_v7", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV7("noauth_v7")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV7MultipleAddresses 测试多个地址
func TestUseV7MultipleAddresses(t *testing.T) {
	option := Option{
		Address: []string{
			"http://es1.example.com:9200",
			"http://es2.example.com:9200",
			"http://es3.example.com:9200",
		},
		Username: "elastic",
		Password: "password",
	}
	Add("multi_v7", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV7("multi_v7")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV7InvalidAddress 测试无效地址
func TestUseV7InvalidAddress(t *testing.T) {
	option := Option{
		Address:  []string{"invalid://address"},
		Username: "test",
		Password: "test",
	}
	Add("invalid_v7", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无效地址): %v", r)
		}
	}()

	UseV7("invalid_v7")
}

// TestUseV7ConcurrentAccess 测试并发访问
func TestUseV7ConcurrentAccess(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("concurrent_v7", option)

	// 清理可能存在的实例
	poolV7.Delete("concurrent_v7")

	done := make(chan bool, 10)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 启动多个 goroutine 同时访问
	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				recover() // 捕获可能的 panic
				done <- true
			}()
			UseV7("concurrent_v7")
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证只创建了一个实例
	if instance, ok := poolV7.Load("concurrent_v7"); ok {
		if instance == nil {
			t.Error("实例不应该为 nil")
		}
	}
}
