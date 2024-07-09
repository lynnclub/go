package rabbit

import (
	"sync"

	"github.com/wagslane/go-rabbitmq"
)

var (
	pool    = &sync.Map{}             //实例池
	mutex   sync.Mutex                //互斥锁
	options = make(map[string]Option) //配置池
)

type Option struct {
	DSN string `json:"dsn"` //数据源
}

func Add(name string, option Option) {
	if option.DSN == "" {
		panic("Option dsn empty " + name)
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	option := Option{
		DSN: setting["dsn"].(string),
	}

	Add(name, option)
}

func AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		AddMap(name, setting.(map[string]interface{}))
	}
}

func Use(name string) *rabbitmq.Conn {
	if name == "" {
		name = "default"
	}

	if instance, ok := pool.Load(name); ok {
		return instance.(*rabbitmq.Conn)
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*rabbitmq.Conn)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	conn, err := rabbitmq.NewConn(option.DSN)
	if err != nil {
		panic("Failed to connect rabbitmq " + name + " err: " + err.Error())
	}

	pool.Store(name, conn)
	return conn
}
