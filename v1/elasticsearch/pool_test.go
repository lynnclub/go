package elasticsearch

import (
	"testing"
)

// TestAdd 测试添加配置
func TestAdd(t *testing.T) {
	option := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "elastic",
		Password: "password",
	}

	Add("test_add", option)

	if savedOption, ok := options["test_add"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if len(savedOption.Address) != 1 {
			t.Errorf("期望1个地址，实际为%d个", len(savedOption.Address))
		}
		if savedOption.Address[0] != "http://localhost:9200" {
			t.Errorf("期望地址为http://localhost:9200，实际为%s", savedOption.Address[0])
		}
		if savedOption.Username != "elastic" {
			t.Errorf("期望用户名为elastic，实际为%s", savedOption.Username)
		}
		if savedOption.Password != "password" {
			t.Errorf("期望密码为password，实际为%s", savedOption.Password)
		}
	}
}

// TestAddPanic 测试添加空地址配置应该panic
func TestAddPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Add应该在地址为空时panic")
		}
	}()

	option := Option{
		Address: []string{},
	}

	Add("test_panic", option)
}

// TestAddMap 测试从map添加配置
func TestAddMap(t *testing.T) {
	setting := map[string]interface{}{
		"address":  []interface{}{"http://localhost:9200", "http://localhost:9201"},
		"username": "elastic",
		"password": "password",
	}

	AddMap("test_map", setting)

	if savedOption, ok := options["test_map"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if len(savedOption.Address) != 2 {
			t.Errorf("期望2个地址，实际为%d个", len(savedOption.Address))
		}
		if savedOption.Address[0] != "http://localhost:9200" {
			t.Errorf("期望第一个地址为http://localhost:9200，实际为%s", savedOption.Address[0])
		}
		if savedOption.Address[1] != "http://localhost:9201" {
			t.Errorf("期望第二个地址为http://localhost:9201，实际为%s", savedOption.Address[1])
		}
		if savedOption.Username != "elastic" {
			t.Errorf("期望用户名为elastic，实际为%s", savedOption.Username)
		}
		if savedOption.Password != "password" {
			t.Errorf("期望密码为password，实际为%s", savedOption.Password)
		}
	}
}

// TestAddMapWithoutAuth 测试从map添加不带认证的配置
func TestAddMapWithoutAuth(t *testing.T) {
	setting := map[string]interface{}{
		"address": []interface{}{"http://localhost:9200"},
	}

	AddMap("test_map_no_auth", setting)

	if savedOption, ok := options["test_map_no_auth"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if len(savedOption.Address) != 1 {
			t.Errorf("期望1个地址，实际为%d个", len(savedOption.Address))
		}
		if savedOption.Username != "" {
			t.Errorf("期望用户名为空，实际为%s", savedOption.Username)
		}
		if savedOption.Password != "" {
			t.Errorf("期望密码为空，实际为%s", savedOption.Password)
		}
	}
}

// TestAddMapBatch 测试批量添加配置
func TestAddMapBatch(t *testing.T) {
	batch := map[string]interface{}{
		"es1": map[string]interface{}{
			"address":  []interface{}{"http://localhost:9200"},
			"username": "user1",
			"password": "pass1",
		},
		"es2": map[string]interface{}{
			"address":  []interface{}{"http://localhost:9201"},
			"username": "user2",
			"password": "pass2",
		},
	}

	AddMapBatch(batch)

	// 验证es1
	if savedOption, ok := options["es1"]; !ok {
		t.Error("es1配置未成功保存")
	} else {
		if savedOption.Username != "user1" {
			t.Errorf("期望es1用户名为user1，实际为%s", savedOption.Username)
		}
	}

	// 验证es2
	if savedOption, ok := options["es2"]; !ok {
		t.Error("es2配置未成功保存")
	} else {
		if savedOption.Username != "user2" {
			t.Errorf("期望es2用户名为user2，实际为%s", savedOption.Username)
		}
	}
}

// TestAddMapBatchEmpty 测试批量添加空配置
func TestAddMapBatchEmpty(t *testing.T) {
	batch := map[string]interface{}{}

	// 不应该panic
	AddMapBatch(batch)
}

// TestAddMultipleAddresses 测试添加多个地址
func TestAddMultipleAddresses(t *testing.T) {
	option := Option{
		Address: []string{
			"http://es1.example.com:9200",
			"http://es2.example.com:9200",
			"http://es3.example.com:9200",
		},
		Username: "elastic",
		Password: "password",
	}

	Add("test_multi", option)

	if savedOption, ok := options["test_multi"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if len(savedOption.Address) != 3 {
			t.Errorf("期望3个地址，实际为%d个", len(savedOption.Address))
		}
	}
}

// TestOptionOverwrite 测试配置覆盖
func TestOptionOverwrite(t *testing.T) {
	option1 := Option{
		Address:  []string{"http://localhost:9200"},
		Username: "user1",
		Password: "pass1",
	}

	Add("test_overwrite", option1)

	option2 := Option{
		Address:  []string{"http://localhost:9201"},
		Username: "user2",
		Password: "pass2",
	}

	Add("test_overwrite", option2)

	// 验证配置被覆盖
	if savedOption, ok := options["test_overwrite"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.Address[0] != "http://localhost:9201" {
			t.Errorf("期望地址为http://localhost:9201，实际为%s", savedOption.Address[0])
		}
		if savedOption.Username != "user2" {
			t.Errorf("期望用户名为user2，实际为%s", savedOption.Username)
		}
	}
}
