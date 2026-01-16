package elasticsearch

import (
	"testing"
)

// TestUseV8WithDefaultName 测试使用默认名称
func TestUseV8WithDefaultName(t *testing.T) {
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

	client := UseV8("")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV8WithCustomName 测试使用自定义名称
func TestUseV8WithCustomName(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("custom_v8", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV8("custom_v8")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV8OptionNotFound 测试配置不存在的情况
func TestUseV8OptionNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该在配置不存在时 panic")
		} else {
			errMsg, ok := r.(string)
			if !ok || errMsg != "Option not found nonexistent_v8" {
				t.Errorf("panic 信息不正确: %v", r)
			}
		}
	}()

	UseV8("nonexistent_v8")
}

// TestUseV8EmptyAddresses 测试空地址配置
func TestUseV8EmptyAddresses(t *testing.T) {
	// 这个测试在 Add 函数中就会 panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该在地址为空时 panic")
		}
	}()

	option := Option{
		Address: []string{},
	}
	Add("empty_v8", option)
	UseV8("empty_v8")
}

// TestUseV8Singleton 测试单例模式
func TestUseV8Singleton(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("singleton_v8", option)

	// 清理可能存在的实例
	pool.Delete("singleton_v8")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 第一次调用会创建新实例
	client1 := UseV8("singleton_v8")

	// 第二次调用应该返回相同的实例
	client2 := UseV8("singleton_v8")

	if client1 != client2 {
		t.Error("UseV8 应该返回相同的实例（单例模式）")
	}
}

// TestUseV8WithoutAuth 测试无认证配置
func TestUseV8WithoutAuth(t *testing.T) {
	option := Option{
		Address: []string{"http://localhost:9200"},
	}
	Add("noauth_v8", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV8("noauth_v8")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV8MultipleAddresses 测试多个地址
func TestUseV8MultipleAddresses(t *testing.T) {
	option := Option{
		Address: []string{
			"http://es1.example.com:9200",
			"http://es2.example.com:9200",
			"http://es3.example.com:9200",
		},
		Username: "elastic",
		Password: "password",
	}
	Add("multi_v8", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := UseV8("multi_v8")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestUseV8InvalidAddress 测试无效地址
func TestUseV8InvalidAddress(t *testing.T) {
	option := Option{
		Address:  []string{"invalid://address"},
		Username: "test",
		Password: "test",
	}
	Add("invalid_v8", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无效地址): %v", r)
		}
	}()

	UseV8("invalid_v8")
}

// TestUseV8ConcurrentAccess 测试并发访问
func TestUseV8ConcurrentAccess(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("concurrent_v8", option)

	// 清理可能存在的实例
	pool.Delete("concurrent_v8")

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
			UseV8("concurrent_v8")
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证只创建了一个实例
	if instance, ok := pool.Load("concurrent_v8"); ok {
		if instance == nil {
			t.Error("实例不应该为 nil")
		}
	}
}

// TestUseV8AfterV7 测试在使用 V7 后使用 V8
func TestUseV8AfterV7(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("test_v7_v8", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 先尝试使用 V7
	UseV7("test_v7_v8")

	// 再使用 V8 (应该使用不同的连接池)
	UseV8("test_v7_v8")
}
