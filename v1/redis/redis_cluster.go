package redis

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	poolCluster  = &sync.Map{} //实例池
	mutexCluster sync.Mutex    //互斥锁
)

// Cluster 使用集群
func Cluster(name string) *redis.ClusterClient {
	if instance, ok := poolCluster.Load(name); ok {
		return instance.(*redis.ClusterClient)
	} else {
		mutexCluster.Lock()
		defer mutexCluster.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*redis.ClusterClient)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	newClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    option.Address,
		Password: option.Password,
		PoolSize: option.PoolSize,
	})

	_, err := newClient.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect redis cluster " + name + " err: " + err.Error())
	}

	poolCluster.Store(name, newClient)
	return newClient
}
