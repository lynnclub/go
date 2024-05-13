package redis

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	poolCluster = &sync.Map{} //实例池
)

// Cluster 使用集群
func Cluster(name string) *redis.ClusterClient {
	if instance, ok := poolCluster.Load(name); ok {
		return instance.(*redis.ClusterClient)
	} else {
		var mutex sync.Mutex
		mutex.Lock()
		defer mutex.Unlock()
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	newCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    option.Address,
		Password: option.Password,
		PoolSize: option.PoolSize,
	})

	poolCluster.Store(name, newCluster)
	return newCluster
}
