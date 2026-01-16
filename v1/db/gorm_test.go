package db

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestAdd 测试添加配置
func TestAdd(t *testing.T) {
	option := Option{
		DSN:    ":memory:",
		Driver: "sqlite",
	}
	Add("test_add", option)

	if _, ok := options["test_add"]; !ok {
		t.Error("添加配置失败")
	}

	// 验证默认值
	opt := options["test_add"]
	if opt.Driver != "sqlite" {
		t.Errorf("期望Driver为sqlite，实际为%s", opt.Driver)
	}
	if opt.MaxOpenConn != 100 {
		t.Errorf("期望MaxOpenConn为100，实际为%d", opt.MaxOpenConn)
	}
	if opt.MaxIdleConn != 10 {
		t.Errorf("期望MaxIdleConn为10，实际为%d", opt.MaxIdleConn)
	}
	if opt.MaxIdleTime != 600 {
		t.Errorf("期望MaxIdleTime为600，实际为%d", opt.MaxIdleTime)
	}
	if opt.LogLevel != 3 {
		t.Errorf("期望LogLevel为3，实际为%d", opt.LogLevel)
	}
	if opt.SlowThreshold != 1 {
		t.Errorf("期望SlowThreshold为1，实际为%d", opt.SlowThreshold)
	}
}

// TestAddWithCustomValues 测试添加自定义配置
func TestAddWithCustomValues(t *testing.T) {
	option := Option{
		DSN:           ":memory:",
		Driver:        "sqlite",
		MaxOpenConn:   50,
		MaxIdleConn:   5,
		MaxIdleTime:   300,
		LogLevel:      4,
		SlowThreshold: 2,
	}
	Add("test_custom", option)

	opt := options["test_custom"]
	if opt.MaxOpenConn != 50 {
		t.Errorf("期望MaxOpenConn为50，实际为%d", opt.MaxOpenConn)
	}
	if opt.MaxIdleConn != 5 {
		t.Errorf("期望MaxIdleConn为5，实际为%d", opt.MaxIdleConn)
	}
	if opt.MaxIdleTime != 300 {
		t.Errorf("期望MaxIdleTime为300，实际为%d", opt.MaxIdleTime)
	}
	if opt.LogLevel != 4 {
		t.Errorf("期望LogLevel为4，实际为%d", opt.LogLevel)
	}
	if opt.SlowThreshold != 2 {
		t.Errorf("期望SlowThreshold为2，实际为%d", opt.SlowThreshold)
	}
}

// TestAddMap 测试从map添加配置
func TestAddMap(t *testing.T) {
	setting := map[string]interface{}{
		"dsn":            ":memory:",
		"driver":         "sqlite",
		"max_open_conn":  80,
		"max_idle_conn":  8,
		"max_idle_time":  400,
		"log_level":      2,
		"slow_threshold": 3,
	}
	AddMap("test_add_map", setting)

	opt, ok := options["test_add_map"]
	if !ok {
		t.Error("从map添加配置失败")
	}

	if opt.DSN != ":memory:" {
		t.Errorf("期望DSN为:memory:，实际为%s", opt.DSN)
	}
	if opt.Driver != "sqlite" {
		t.Errorf("期望Driver为sqlite，实际为%s", opt.Driver)
	}
	if opt.MaxOpenConn != 80 {
		t.Errorf("期望MaxOpenConn为80，实际为%d", opt.MaxOpenConn)
	}
}

// TestAddMapBatch 测试批量添加配置
func TestAddMapBatch(t *testing.T) {
	batch := map[string]interface{}{
		"batch1": map[string]interface{}{
			"dsn":    ":memory:",
			"driver": "sqlite",
		},
		"batch2": map[string]interface{}{
			"dsn":    ":memory:",
			"driver": "sqlite",
		},
	}
	AddMapBatch(batch)

	if _, ok := options["batch1"]; !ok {
		t.Error("批量添加配置batch1失败")
	}
	if _, ok := options["batch2"]; !ok {
		t.Error("批量添加配置batch2失败")
	}
}

// TestUse 测试使用数据库连接
func TestUse(t *testing.T) {
	// 使用SQLite内存数据库
	option := Option{
		DSN:    ":memory:",
		Driver: "sqlite",
	}
	Add("test_use", option)

	db := Use("test_use")
	if db == nil {
		t.Error("获取数据库连接失败")
	}

	// 测试连接复用
	if db != Use("test_use") {
		t.Error("连接未能复用")
	}

	// 测试数据库操作
	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("获取sql.DB失败: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		t.Errorf("数据库Ping失败: %v", err)
	}

	// 测试连接池配置
	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 100 {
		t.Errorf("期望MaxOpenConnections为100，实际为%d", stats.MaxOpenConnections)
	}
}

// TestUseDefault 测试使用默认连接
func TestUseDefault(t *testing.T) {
	option := Option{
		DSN:    ":memory:",
		Driver: "sqlite",
	}
	Add("default", option)

	// 空字符串应该使用default
	db := Use("")
	if db == nil {
		t.Error("获取默认连接失败")
	}

	if db != Use("default") {
		t.Error("空字符串未映射到default连接")
	}
}

// TestUseConcurrent 测试并发使用数据库连接
func TestUseConcurrent(t *testing.T) {
	option := Option{
		DSN:    ":memory:",
		Driver: "sqlite",
	}
	Add("test_concurrent", option)

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for loop := 0; loop < 10; loop++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			db := Use("test_concurrent")
			if db == nil {
				errors <- fmt.Errorf("goroutine %d: 获取连接失败", i)
				return
			}

			// 测试数据库操作
			sqlDB, err := db.DB()
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: 获取sql.DB失败: %v", i, err)
				return
			}

			if err = sqlDB.Ping(); err != nil {
				errors <- fmt.Errorf("goroutine %d: Ping失败: %v", i, err)
			}
		}(loop)
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	for err := range errors {
		t.Error(err)
	}
}

// TestDrivers 测试不同的数据库驱动
func TestDrivers(t *testing.T) {
	drivers := []string{"mysql", "postgres", "sqlite", "sqlserver", "clickhouse"}

	for _, driver := range drivers {
		t.Run(driver, func(t *testing.T) {
			// 只测试配置，不实际连接（除了sqlite）
			var dsn string
			switch driver {
			case "sqlite":
				dsn = ":memory:"
			case "mysql":
				dsn = "user:pass@tcp(localhost:3306)/dbname"
			case "postgres":
				dsn = "host=localhost user=user password=pass dbname=db"
			case "sqlserver":
				dsn = "sqlserver://user:pass@localhost:1433?database=db"
			case "clickhouse":
				dsn = "tcp://localhost:9000?database=default"
			}

			option := Option{
				DSN:    dsn,
				Driver: driver,
			}
			Add("test_driver_"+driver, option)

			opt := options["test_driver_"+driver]
			if opt.Driver != driver {
				t.Errorf("期望Driver为%s，实际为%s", driver, opt.Driver)
			}

			// 只对sqlite进行实际连接测试
			if driver == "sqlite" {
				db := Use("test_driver_" + driver)
				if db == nil {
					t.Error("SQLite连接失败")
				}
			}
		})
	}
}

// TestUsePanic 测试使用不存在的配置
func TestUsePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望panic但没有发生")
		}
	}()

	Use("non_existent_db_config")
}

// TestAddPanic 测试添加空DSN配置
func TestAddPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望panic但没有发生")
		}
	}()

	option := Option{
		DSN: "", // 空DSN
	}
	Add("test_panic", option)
}

// TestConnectionPoolSettings 测试连接池设置
func TestConnectionPoolSettings(t *testing.T) {
	option := Option{
		DSN:         ":memory:",
		Driver:      "sqlite",
		MaxOpenConn: 50,
		MaxIdleConn: 5,
		MaxIdleTime: 300,
	}
	Add("test_pool", option)

	db := Use("test_pool")
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取sql.DB失败: %v", err)
	}

	// 验证连接池设置
	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 50 {
		t.Errorf("期望MaxOpenConnections为50，实际为%d", stats.MaxOpenConnections)
	}

	// 设置后立即检查
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxIdleTime(300 * time.Second)

	// SQLite的Stats可能不完全准确，但至少不应该报错
}

// TestLogLevels 测试不同的日志级别
func TestLogLevels(t *testing.T) {
	levels := []int{1, 2, 3, 4} // Silent, Error, Warn, Info

	for _, level := range levels {
		t.Run(fmt.Sprintf("LogLevel_%d", level), func(t *testing.T) {
			option := Option{
				DSN:      ":memory:",
				Driver:   "sqlite",
				LogLevel: level,
			}
			Add(fmt.Sprintf("test_log_%d", level), option)

			opt := options[fmt.Sprintf("test_log_%d", level)]
			if opt.LogLevel != level {
				t.Errorf("期望LogLevel为%d，实际为%d", level, opt.LogLevel)
			}

			// 实际连接测试
			db := Use(fmt.Sprintf("test_log_%d", level))
			if db == nil {
				t.Error("连接失败")
			}
		})
	}
}
