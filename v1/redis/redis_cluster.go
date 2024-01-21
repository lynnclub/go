package redis

import (
	"github.com/go-redis/redis/v8"
)

var (
	poolCluster = make(map[string]*redis.ClusterClient)
)

// Cluster 使用集群
func Cluster(name string) *redis.ClusterClient {
	if instance, ok := poolCluster[name]; ok {
		return instance
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	poolCluster[name] = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    option.Address,
		Password: option.Password,
		PoolSize: option.PoolSize,
	})

	return poolCluster[name]
}
