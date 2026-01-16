package mongodb

import (
	"testing"
)

// TestAdd 测试添加配置
func TestAdd(t *testing.T) {
	option := Option{
		DSN: "mongodb://localhost:27017",
	}

	Add("test_add", option)

	if savedOption, ok := options["test_add"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://localhost:27017" {
			t.Errorf("期望DSN为mongodb://localhost:27017，实际为%s", savedOption.DSN)
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
		"dsn": "mongodb://user:pass@localhost:27017/testdb",
	}

	AddMap("test_map", setting)

	if savedOption, ok := options["test_map"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://user:pass@localhost:27017/testdb" {
			t.Errorf("期望DSN为mongodb://user:pass@localhost:27017/testdb，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddMapWithAuth 测试带认证的配置
func TestAddMapWithAuth(t *testing.T) {
	setting := map[string]interface{}{
		"dsn": "mongodb://admin:password@localhost:27017/admin?authSource=admin",
	}

	AddMap("test_auth", setting)

	if savedOption, ok := options["test_auth"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://admin:password@localhost:27017/admin?authSource=admin"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddMapBatch 测试批量添加配置
func TestAddMapBatch(t *testing.T) {
	batch := map[string]interface{}{
		"mongo1": map[string]interface{}{
			"dsn": "mongodb://localhost:27017/db1",
		},
		"mongo2": map[string]interface{}{
			"dsn": "mongodb://localhost:27017/db2",
		},
	}

	AddMapBatch(batch)

	// 验证mongo1
	if savedOption, ok := options["mongo1"]; !ok {
		t.Error("mongo1配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://localhost:27017/db1" {
			t.Errorf("期望mongo1 DSN为mongodb://localhost:27017/db1，实际为%s", savedOption.DSN)
		}
	}

	// 验证mongo2
	if savedOption, ok := options["mongo2"]; !ok {
		t.Error("mongo2配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://localhost:27017/db2" {
			t.Errorf("期望mongo2 DSN为mongodb://localhost:27017/db2，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddMapBatchEmpty 测试批量添加空配置
func TestAddMapBatchEmpty(t *testing.T) {
	batch := map[string]interface{}{}

	// 不应该panic
	AddMapBatch(batch)
}

// TestAddReplicaSet 测试添加副本集配置
func TestAddReplicaSet(t *testing.T) {
	option := Option{
		DSN: "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myReplicaSet",
	}

	Add("test_replica", option)

	if savedOption, ok := options["test_replica"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myReplicaSet"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithOptions 测试带选项的配置
func TestAddWithOptions(t *testing.T) {
	option := Option{
		DSN: "mongodb://localhost:27017/testdb?maxPoolSize=100&minPoolSize=10",
	}

	Add("test_options", option)

	if savedOption, ok := options["test_options"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://localhost:27017/testdb?maxPoolSize=100&minPoolSize=10"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestOptionOverwrite 测试配置覆盖
func TestOptionOverwrite(t *testing.T) {
	option1 := Option{
		DSN: "mongodb://localhost:27017/db1",
	}

	Add("test_overwrite", option1)

	option2 := Option{
		DSN: "mongodb://localhost:27017/db2",
	}

	Add("test_overwrite", option2)

	// 验证配置被覆盖
	if savedOption, ok := options["test_overwrite"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://localhost:27017/db2" {
			t.Errorf("期望DSN为mongodb://localhost:27017/db2，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddSRVConnection 测试SRV连接字符串
func TestAddSRVConnection(t *testing.T) {
	option := Option{
		DSN: "mongodb+srv://user:pass@cluster.mongodb.net/testdb",
	}

	Add("test_srv", option)

	if savedOption, ok := options["test_srv"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb+srv://user:pass@cluster.mongodb.net/testdb"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddLocalhost 测试本地连接
func TestAddLocalhost(t *testing.T) {
	option := Option{
		DSN: "mongodb://127.0.0.1:27017",
	}

	Add("test_localhost", option)

	if savedOption, ok := options["test_localhost"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://127.0.0.1:27017" {
			t.Errorf("期望DSN为mongodb://127.0.0.1:27017，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddWithDatabase 测试指定数据库的配置
func TestAddWithDatabase(t *testing.T) {
	option := Option{
		DSN: "mongodb://localhost:27017/myapp",
	}

	Add("test_db", option)

	if savedOption, ok := options["test_db"]; !ok {
		t.Error("配置未成功保存")
	} else {
		if savedOption.DSN != "mongodb://localhost:27017/myapp" {
			t.Errorf("期望DSN为mongodb://localhost:27017/myapp，实际为%s", savedOption.DSN)
		}
	}
}

// TestAddMultipleHosts 测试多主机配置
func TestAddMultipleHosts(t *testing.T) {
	option := Option{
		DSN: "mongodb://host1:27017,host2:27018,host3:27019/testdb",
	}

	Add("test_multi_host", option)

	if savedOption, ok := options["test_multi_host"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://host1:27017,host2:27018,host3:27019/testdb"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithReadPreference 测试带读偏好的配置
func TestAddWithReadPreference(t *testing.T) {
	option := Option{
		DSN: "mongodb://localhost:27017/testdb?readPreference=secondary",
	}

	Add("test_read_pref", option)

	if savedOption, ok := options["test_read_pref"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://localhost:27017/testdb?readPreference=secondary"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithWriteConcern 测试带写关注的配置
func TestAddWithWriteConcern(t *testing.T) {
	option := Option{
		DSN: "mongodb://localhost:27017/testdb?w=majority&wtimeoutMS=5000",
	}

	Add("test_write_concern", option)

	if savedOption, ok := options["test_write_concern"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://localhost:27017/testdb?w=majority&wtimeoutMS=5000"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}

// TestAddWithSSL 测试带SSL的配置
func TestAddWithSSL(t *testing.T) {
	option := Option{
		DSN: "mongodb://localhost:27017/testdb?ssl=true",
	}

	Add("test_ssl", option)

	if savedOption, ok := options["test_ssl"]; !ok {
		t.Error("配置未成功保存")
	} else {
		expected := "mongodb://localhost:27017/testdb?ssl=true"
		if savedOption.DSN != expected {
			t.Errorf("期望DSN为%s，实际为%s", expected, savedOption.DSN)
		}
	}
}
