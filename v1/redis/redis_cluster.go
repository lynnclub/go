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
		mutex.Lock()
		defer mutex.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*redis.ClusterClient)
		}
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

	_, err := newCluster.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect redis cluster " + name + " err: " + err.Error())
	}

	poolCluster.Store(name, newCluster)
	return newCluster
}
