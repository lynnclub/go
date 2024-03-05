package db

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	pool    = &sync.Map{}             //实例池
	options = make(map[string]Option) //配置池
	Default *gorm.DB                  //默认数据库
)

type Option struct {
	DSN           string `json:"dsn"`            //数据源
	Driver        string `json:"driver"`         //驱动，mysql、postgres、sqlite、sqlserver、clickhouse，默认mysql
	MaxOpenConn   int    `json:"max_open_conn"`  //打开连接的最大数量，默认100
	MaxIdleConn   int    `json:"max_idle_conn"`  //空闲连接的最大数量，默认10
	LogLevel      int    `json:"log_level"`      //Silent 1、Error 2、Warn 3、Info 4，默认3
	SlowThreshold int    `json:"slow_threshold"` //慢SQL阈值，单位秒，默认1
}

func Add(name string, option Option) {
	if option.DSN == "" {
		panic("Option dsn empty " + name)
	}

	// 默认值
	if option.Driver == "" {
		option.Driver = "mysql"
	}
	if option.MaxOpenConn == 0 {
		option.MaxOpenConn = 100
	}
	if option.MaxIdleConn == 0 {
		option.MaxIdleConn = 10
	}
	if option.LogLevel == 0 {
		option.LogLevel = 3
	}
	if option.SlowThreshold == 0 {
		option.SlowThreshold = 1
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	option := Option{
		DSN: setting["dsn"].(string),
	}

	if driver, ok := setting["driver"]; ok {
		option.Driver = driver.(string)
	}
	if maxOpenConn, ok := setting["max_open_conn"]; ok {
		option.MaxOpenConn = maxOpenConn.(int)
	}
	if maxIdleConn, ok := setting["max_idle_conn"]; ok {
		option.MaxIdleConn = maxIdleConn.(int)
	}
	if logLevel, ok := setting["log_level"]; ok {
		option.LogLevel = logLevel.(int)
	}
	if slowThreshold, ok := setting["slow_threshold"]; ok {
		option.SlowThreshold = slowThreshold.(int)
	}

	Add(name, option)
}

func AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		AddMap(name, setting.(map[string]interface{}))
	}
}

// Use 使用
func Use(name string) *gorm.DB {
	if name == "" {
		name = "default"
	}

	if instance, ok := pool.Load(name); ok {
		return instance.(*gorm.DB)
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	// 驱动
	var dialect gorm.Dialector
	switch option.Driver {
	case "mysql":
		dialect = mysql.Open(option.DSN)
	case "postgres":
		dialect = postgres.Open(option.DSN)
	case "sqlite":
		dialect = sqlite.Open(option.DSN)
	case "sqlserver":
		dialect = sqlserver.Open(option.DSN)
	case "clickhouse":
		dialect = clickhouse.Open(option.DSN)
	default:
		panic("Driver not support " + option.Driver)
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Duration(option.SlowThreshold) * time.Second, // 慢SQL阈值
			LogLevel:                  logger.LogLevel(option.LogLevel),                  // 日志级别
			IgnoreRecordNotFoundError: true,                                              // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                                             // 彩色打印
		},
	)

	newGorm, err := gorm.Open(dialect, &gorm.Config{Logger: gormLogger})
	if err != nil {
		panic("Failed to connect database " + name + " err: " + err.Error())
	}

	sqlDB, _ := newGorm.DB()
	// 打开连接的最大数量
	sqlDB.SetMaxOpenConns(option.MaxOpenConn)
	// 空闲连接的最大数量
	sqlDB.SetMaxIdleConns(option.MaxIdleConn)

	if name == "default" {
		Default = newGorm
	}

	pool.Store(name, newGorm)
	return newGorm
}
