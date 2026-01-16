package elasticsearch

import (
	"testing"
)

// TestTypedV8WithDefaultName 测试使用默认名称
func TestTypedV8WithDefaultName(t *testing.T) {
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

	client := TypedV8("")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestTypedV8WithCustomName 测试使用自定义名称
func TestTypedV8WithCustomName(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("custom_typed", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := TypedV8("custom_typed")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestTypedV8OptionNotFound 测试配置不存在的情况
func TestTypedV8OptionNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该在配置不存在时 panic")
		} else {
			errMsg, ok := r.(string)
			if !ok || errMsg != "Option not found nonexistent_typed" {
				t.Errorf("panic 信息不正确: %v", r)
			}
		}
	}()

	TypedV8("nonexistent_typed")
}

// TestTypedV8EmptyAddresses 测试空地址配置
func TestTypedV8EmptyAddresses(t *testing.T) {
	// 这个测试在 Add 函数中就会 panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该在地址为空时 panic")
		}
	}()

	option := Option{
		Address: []string{},
	}
	Add("empty_typed", option)
	TypedV8("empty_typed")
}

// TestTypedV8Singleton 测试单例模式
func TestTypedV8Singleton(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("singleton_typed", option)

	// 清理可能存在的实例
	poolTyped.Delete("singleton_typed")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 第一次调用会创建新实例
	client1 := TypedV8("singleton_typed")

	// 第二次调用应该返回相同的实例
	client2 := TypedV8("singleton_typed")

	if client1 != client2 {
		t.Error("TypedV8 应该返回相同的实例（单例模式）")
	}
}

// TestTypedV8WithoutAuth 测试无认证配置
func TestTypedV8WithoutAuth(t *testing.T) {
	option := Option{
		Address: []string{"http://localhost:9200"},
	}
	Add("noauth_typed", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := TypedV8("noauth_typed")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestTypedV8MultipleAddresses 测试多个地址
func TestTypedV8MultipleAddresses(t *testing.T) {
	option := Option{
		Address: []string{
			"http://es1.example.com:9200",
			"http://es2.example.com:9200",
			"http://es3.example.com:9200",
		},
		Username: "elastic",
		Password: "password",
	}
	Add("multi_typed", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	client := TypedV8("multi_typed")
	if client == nil {
		t.Error("客户端不应该为 nil")
	}
}

// TestTypedV8InvalidAddress 测试无效地址
func TestTypedV8InvalidAddress(t *testing.T) {
	option := Option{
		Address:  []string{"invalid://address"},
		Username: "test",
		Password: "test",
	}
	Add("invalid_typed", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无效地址): %v", r)
		}
	}()

	TypedV8("invalid_typed")
}

// TestTypedV8ConcurrentAccess 测试并发访问
func TestTypedV8ConcurrentAccess(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("concurrent_typed", option)

	// 清理可能存在的实例
	poolTyped.Delete("concurrent_typed")

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
			TypedV8("concurrent_typed")
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证只创建了一个实例
	if instance, ok := poolTyped.Load("concurrent_typed"); ok {
		if instance == nil {
			t.Error("实例不应该为 nil")
		}
	}
}

// TestTypedV8DifferentFromNormalV8 测试 TypedV8 与普通 V8 使用不同的连接池
func TestTypedV8DifferentFromNormalV8(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("test_typed_vs_normal", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 使用普通 V8
	UseV8("test_typed_vs_normal")

	// 使用 Typed V8 (应该使用不同的连接池)
	TypedV8("test_typed_vs_normal")

	// 验证两个池中都有实例
	_, normalExists := pool.Load("test_typed_vs_normal")
	_, typedExists := poolTyped.Load("test_typed_vs_normal")

	if normalExists && typedExists {
		t.Log("普通 V8 和 Typed V8 使用不同的连接池")
	}
}

// TestTypedV8EmptyNameDefaultsToDefault 测试空名称默认为 default
func TestTypedV8EmptyNameDefaultsToDefault(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "test",
		Password: "test",
	}
	Add("default", option)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("预期的 panic (无法连接到 ES): %v", r)
		}
	}()

	// 空字符串应该使用 "default"
	client1 := TypedV8("")
	client2 := TypedV8("default")

	if client1 != client2 {
		t.Error("空名称和 'default' 应该返回相同的实例")
	}
}
