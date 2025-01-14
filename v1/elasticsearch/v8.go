package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

// UseV8 使用
func UseV8(name string) *elasticsearch.Client {
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
	if err == nil {
		fmt.Println("Connected to elasticsearch", name, info)
		info.Body.Close()
	} else {
		panic("Failed to connect elasticsearch " + name + " err: " + err.Error())
	}

	pool.Store(name, newClient)
	return newClient
}
