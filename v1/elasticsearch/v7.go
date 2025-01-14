package elasticsearch

import (
	"fmt"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
)

var (
	poolV7  = &sync.Map{} //实例池
	mutexV7 sync.Mutex    //互斥锁
)

// UseV7 使用V7
func UseV7(name string) *elasticsearch.Client {
	if name == "" {
		name = "default"
	}

	if instance, ok := poolV7.Load(name); ok {
		return instance.(*elasticsearch.Client)
	} else {
		mutexV7.Lock()
		defer mutexV7.Unlock()
		if instance, ok = poolV7.Load(name); ok {
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
