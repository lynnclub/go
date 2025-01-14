package elasticsearch

import (
	"context"
	"fmt"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	poolTyped  = &sync.Map{} //实例池
	mutexTyped sync.Mutex    //互斥锁
)

// TypedV8 使用V8 TypedAPI
func TypedV8(name string) *elasticsearch.TypedClient {
	if name == "" {
		name = "default"
	}

	if instance, ok := poolTyped.Load(name); ok {
		return instance.(*elasticsearch.TypedClient)
	} else {
		mutexTyped.Lock()
		defer mutexTyped.Unlock()
		if instance, ok = poolTyped.Load(name); ok {
			return instance.(*elasticsearch.TypedClient)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	newClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: option.Address,
		Username:  option.Username,
		Password:  option.Password,
	})
	if err != nil {
		panic("Failed to new elasticsearch " + name + " err: " + err.Error())
	}

	info, err := newClient.Info().Do(context.Background())
	if err == nil {
		fmt.Println("Connected to elasticsearch", name, info)
	} else {
		panic("Failed to connect elasticsearch " + name + " err: " + err.Error())
	}

	pool.Store(name, newClient)
	return newClient
}
