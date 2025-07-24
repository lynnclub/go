package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	pool    = &sync.Map{}             //实例池
	mutex   sync.Mutex                //互斥锁
	options = make(map[string]Option) //配置池
	Ctx     = context.Background()
	Nil     = redis.Nil
)

type Option struct {
	Address         []string      `json:"address"`            //地址，字符串数组
	Password        string        `json:"password"`           //密码，默认空
	DB              int           `json:"db"`                 //db
	PoolSize        int           `json:"pool_size"`          //连接池最大数量，默认100
	MinIdleConns    int           `json:"min_idle_conns"`     //最小空闲连接数，默认0
	MaxIdleConns    int           `json:"max_idle_conns"`     //最大空闲连接数，默认0（无限制）
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"` //连接最大空闲时间，默认5分钟
	MasterName      string        `json:"master_name"`        //Sentinel集群模式，主库名称，默认mymaster
	TLS             bool          `json:"tls"`                //是否启用TLS，默认使用系统根证书
}

func Add(name string, option Option) {
	if len(option.Address) == 0 {
		panic("Option address array empty " + name)
	}

	// 默认值
	if option.PoolSize == 0 {
		option.PoolSize = 100
	}
	if option.MasterName == "" {
		option.MasterName = "mymaster"
	}
	if option.ConnMaxIdleTime <= 0 {
		option.ConnMaxIdleTime = 5 * time.Minute
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	addressStrings := make([]string, 0)
	for _, v := range setting["address"].([]interface{}) {
		addressStrings = append(addressStrings, v.(string))
	}

	option := Option{
		Address: addressStrings,
	}

	if password, ok := setting["password"]; ok {
		option.Password = password.(string)
	}
	if db, ok := setting["db"]; ok {
		option.DB = db.(int)
	}
	if poolSize, ok := setting["pool_size"]; ok {
		option.PoolSize = poolSize.(int)
	}
	if minIdleConns, ok := setting["min_idle_conns"]; ok {
		option.MinIdleConns = minIdleConns.(int)
	}
	if maxIdleConns, ok := setting["max_idle_conns"]; ok {
		option.MaxIdleConns = maxIdleConns.(int)
	}
	if connMaxIdleTime, ok := setting["conn_max_idle_time"]; ok {
		if duration, isString := connMaxIdleTime.(string); isString {
			if parsed, err := time.ParseDuration(duration); err == nil {
				option.ConnMaxIdleTime = parsed
			}
		} else if duration, isInt := connMaxIdleTime.(int); isInt {
			option.ConnMaxIdleTime = time.Duration(duration) * time.Second
		}
	}
	if masterName, ok := setting["master_name"]; ok {
		option.MasterName = masterName.(string)
	}
	if tls, ok := setting["tls"]; ok {
		option.TLS = tls.(bool)
	}

	Add(name, option)
}

func AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		AddMap(name, setting.(map[string]interface{}))
	}
}

// Use 使用
func Use(name string) *redis.Client {
	if name == "" {
		name = "default"
	}

	if instance, ok := pool.Load(name); ok {
		return instance.(*redis.Client)
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*redis.Client)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	clientOptions := &redis.Options{
		Addr:            option.Address[0],
		DB:              option.DB,
		Password:        option.Password,
		PoolSize:        option.PoolSize,
		MinIdleConns:    option.MinIdleConns,
		MaxIdleConns:    option.MaxIdleConns,
		ConnMaxIdleTime: option.ConnMaxIdleTime,
	}
	if option.TLS {
		clientOptions.TLSConfig = &tls.Config{}
	}

	newClient := redis.NewClient(clientOptions)

	info, err := newClient.Ping(Ctx).Result()
	if err == nil {
		fmt.Println("Connected to redis", name, info)
	} else {
		panic("Failed to connect redis " + name + " err: " + err.Error())
	}

	pool.Store(name, newClient)
	return newClient
}
