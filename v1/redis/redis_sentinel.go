package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	poolSentinel  = &sync.Map{} //实例池
	mutexSentinel sync.Mutex    //互斥锁
)

// Sentinel 使用Sentinel集群
func Sentinel(name string) *redis.Client {
	if instance, ok := poolSentinel.Load(name); ok {
		return instance.(*redis.Client)
	} else {
		mutexSentinel.Lock()
		defer mutexSentinel.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*redis.Client)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	newClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    option.MasterName,
		SentinelAddrs: option.Address,
		Password:      option.Password,
		PoolSize:      option.PoolSize,
	})

	_, err := newClient.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect redis sentinel " + name + " err: " + err.Error())
	}

	poolSentinel.Store(name, newClient)
	return newClient
}
