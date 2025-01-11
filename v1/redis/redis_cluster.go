package redis

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
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

	clusterOptions := &redis.ClusterOptions{
		Addrs:    option.Address,
		Password: option.Password,
		PoolSize: option.PoolSize,
	}
	if option.TLS {
		clusterOptions.TLSConfig = &tls.Config{}
	}

	newClient := redis.NewClusterClient(clusterOptions)

	info, err := newClient.Ping(Ctx).Result()
	if err == nil {
		fmt.Println("Connected to redis cluster", name, info)
	} else {
		panic("Failed to connect redis cluster " + name + " err: " + err.Error())
	}

	poolCluster.Store(name, newClient)
	return newClient
}
