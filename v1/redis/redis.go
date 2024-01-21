package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	pool    = make(map[string]*redis.Client)
	options = make(map[string]Option) //配置池
	Default *redis.Client
	Ctx     = context.Background()
	Nil     = redis.Nil
)

type Option struct {
	Address  []string `json:"address"`   //地址，字符串数组
	Password string   `json:"password"`  //密码，默认空
	DB       int      `json:"db"`        //db
	PoolSize int      `json:"pool_size"` //连接池最大数量，默认100
}

func Add(name string, option Option) {
	if len(option.Address) == 0 {
		panic("Option address array empty " + name)
	}

	// 默认值
	if option.PoolSize == 0 {
		option.PoolSize = 100
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	address := setting["address"].([]interface{})
	addressStrings := make([]string, len(address))
	for i, v := range address {
		addressStrings[i] = v.(string)
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

	if instance, ok := pool[name]; ok {
		return instance
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	pool[name] = redis.NewClient(&redis.Options{
		Addr:     option.Address[0],
		DB:       option.DB,
		Password: option.Password,
		PoolSize: option.PoolSize,
	})

	if name == "default" {
		Default = pool[name]
	}

	return pool[name]
}
