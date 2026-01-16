package rabbit

import (
	"testing"
)

// TestAdd 测试添加配置
func TestAdd(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:5672/",
	}

	Add("test_add", option)

	if savedOption, ok := options["test_add"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "amqp://guest:guest@localhost:5672/" {
			t.Errorf("期望DSN为amqp://guest:guest@localhost:5672/，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddPanic 测试添加空DSN配置应该panic
func TestAddPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Add应该在DSN为空时panic")
		}
	}()

	option := Option{
		DSN: "",
	}

	Add("test_panic", option)
}

// TestAddMap 测试从map添加配置
func TestAddMap(t *testing.T) {
	setting := map[string]interface{}{
		"dsn": "amqp://user:pass@localhost:5672/vhost",
	}

	AddMap("test_map", setting)

	if savedOption, ok := options["test_map"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "amqp://user:pass@localhost:5672/vhost" {
			t.Errorf("期望DSN为amqp://user:pass@localhost:5672/vhost，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddMapWithAuth 测试带认证的配置
func TestAddMapWithAuth(t *testing.T) {
	setting := map[string]interface{}{
		"dsn": "amqp://admin:password@localhost:5672/",
	}

	AddMap("test_auth", setting)

	if savedOption, ok := options["test_auth"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://admin:password@localhost:5672/"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddMapBatch 测试批量添加配置
func TestAddMapBatch(t *testing.T) {
	batch := map[string]interface{}{
		"rabbit1": map[string]interface{}{
			"dsn": "amqp://localhost:5672/vhost1",
		},
		"rabbit2": map[string]interface{}{
			"dsn": "amqp://localhost:5672/vhost2",
		},
	}

	AddMapBatch(batch)

	// 验证rabbit1
	if savedOption, ok := options["rabbit1"]; !ok {
		t.Error("rabbit1配置未成功保存")
	} else {
		if savedOption.DSN != "amqp://localhost:5672/vhost1" {
			t.Errorf("期望rabbit1 DSN为amqp://localhost:5672/vhost1，实际为%s", savedOption.DSN)
		}
	}

	// 验证rabbit2
	if savedOption, ok := options["rabbit2"]; !ok {
		t.Error("rabbit2配置未成功保存")
	} else {
		if savedOption.DSN != "amqp://localhost:5672/vhost2" {
			t.Errorf("期望rabbit2 DSN为amqp://localhost:5672/vhost2，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddMapBatchEmpty 测试批量添加空配置
func TestAddMapBatchEmpty(t *testing.T) {
	batch := map[string]interface{}{}

	// 不应该panic
	AddMapBatch(batch)
}

// TestAddWithVhost 测试带虚拟主机的配置
func TestAddWithVhost(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:5672/myvhost",
	}

	Add("test_vhost", option)

	if savedOption, ok := options["test_vhost"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://guest:guest@localhost:5672/myvhost"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithPort 测试指定端口的配置
func TestAddWithPort(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:15672/",
	}

	Add("test_port", option)

	if savedOption, ok := options["test_port"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://guest:guest@localhost:15672/"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestOptionOverwrite 测试配置覆盖
func TestOptionOverwrite(t *testing.T) {
	option1 := Option{
		DSN: "amqp://localhost:5672/vhost1",
	}

	Add("test_overwrite", option1)

	option2 := Option{
		DSN: "amqp://localhost:5672/vhost2",
	}

	Add("test_overwrite", option2)

	// 验证配置被覆盖
	if savedOption, ok := options["test_overwrite"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "amqp://localhost:5672/vhost2" {
			t.Errorf("期望DSN为amqp://localhost:5672/vhost2，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddWithSSL 测试带SSL的配置
func TestAddWithSSL(t *testing.T) {
	option := Option{
		DSN: "amqps://user:pass@localhost:5671/",
	}

	Add("test_ssl", option)

	if savedOption, ok := options["test_ssl"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqps://user:pass@localhost:5671/"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddLocalhost 测试本地连接
func TestAddLocalhost(t *testing.T) {
	option := Option{
		DSN: "amqp://127.0.0.1:5672/",
	}

	Add("test_localhost", option)

	if savedOption, ok := options["test_localhost"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "amqp://127.0.0.1:5672/" {
			t.Errorf("期望DSN为amqp://127.0.0.1:5672/，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddWithParameters 测试带参数的配置
func TestAddWithParameters(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:5672/?heartbeat=10&connection_timeout=5000",
	}

	Add("test_params", option)

	if savedOption, ok := options["test_params"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://guest:guest@localhost:5672/?heartbeat=10&connection_timeout=5000"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddRemoteHost 测试远程主机连接
func TestAddRemoteHost(t *testing.T) {
	option := Option{
		DSN: "amqp://user:pass@rabbitmq.example.com:5672/production",
	}

	Add("test_remote", option)

	if savedOption, ok := options["test_remote"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://user:pass@rabbitmq.example.com:5672/production"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithComplexPassword 测试包含特殊字符的密码
func TestAddWithComplexPassword(t *testing.T) {
	option := Option{
		DSN: "amqp://user:p%40ssw0rd@localhost:5672/",
	}

	Add("test_complex_pass", option)

	if savedOption, ok := options["test_complex_pass"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://user:p%40ssw0rd@localhost:5672/"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddDefaultVhost 测试默认虚拟主机
func TestAddDefaultVhost(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:5672/%2f",
	}

	Add("test_default_vhost", option)

	if savedOption, ok := options["test_default_vhost"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://guest:guest@localhost:5672/%2f"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithChannelMax 测试带通道最大值的配置
func TestAddWithChannelMax(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:5672/?channel_max=100",
	}

	Add("test_channel_max", option)

	if savedOption, ok := options["test_channel_max"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://guest:guest@localhost:5672/?channel_max=100"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithFrameMax 测试带帧最大值的配置
func TestAddWithFrameMax(t *testing.T) {
	option := Option{
		DSN: "amqp://guest:guest@localhost:5672/?frame_max=131072",
	}

	Add("test_frame_max", option)

	if savedOption, ok := options["test_frame_max"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://guest:guest@localhost:5672/?frame_max=131072"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddClusterConnection 测试集群连接配置
func TestAddClusterConnection(t *testing.T) {
	option := Option{
		DSN: "amqp://user:pass@host1:5672,host2:5672,host3:5672/vhost",
	}

	Add("test_cluster", option)

	if savedOption, ok := options["test_cluster"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "amqp://user:pass@host1:5672,host2:5672,host3:5672/vhost"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}
