package elasticsearch

import (
	"fmt"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	pool    = &sync.Map{}             //实例池
	mutex   sync.Mutex                //互斥锁
	options = make(map[string]Option) //配置池
	Default *elasticsearch.Client
)

type Option struct {
	Address  []string `json:"address"`  //地址，字符串数组
	Username string   `json:"username"` //密码，默认空
	Password string   `json:"password"` //密码，默认空
}

func Add(name string, option Option) {
	if len(option.Address) == 0 {
		panic("Option address array empty " + name)
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	addressStrings := make([]string, 0)
	for i, v := range setting["address"].([]interface{}) {
		addressStrings[i] = v.(string)
	}

	option := Option{
		Address: addressStrings,
	}

	if username, ok := setting["username"]; ok {
		option.Username = username.(string)
	}
	if password, ok := setting["password"]; ok {
		option.Password = password.(string)
	}

	Add(name, option)
}

func AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		AddMap(name, setting.(map[string]interface{}))
	}
}

// Use 使用
func Use(name string) *elasticsearch.Client {
	if name == "" {
		name = "default"
	}

	if instance, ok := pool.Load(name); ok {
		return instance.(*elasticsearch.Client)
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*elasticsearch.Client)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	newClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: option.Address,
		Username:  option.Username,
		Password:  option.Password,
	})
	if err != nil {
		panic("Failed to new elasticsearch " + name + " err: " + err.Error())
	}

	info, err := newClient.Info()
	defer info.Body.Close()
	if err == nil {
		fmt.Println("Connected to elasticsearch", name, info)
	} else {
		panic("Failed to connect elasticsearch " + name + " err: " + err.Error())
	}

	if name == "default" {
		Default = newClient
	}

	pool.Store(name, newClient)
	return newClient
}
