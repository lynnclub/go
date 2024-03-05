package mongo

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	pool    = &sync.Map{}             //实例池
	options = make(map[string]Option) //配置池
	Default *mongo.Client             //默认数据库
)

type Option struct {
	DSN string `json:"dsn"` //数据源
}

func Add(name string, option Option) {
	if option.DSN == "" {
		panic("Option dsn empty " + name)
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	option := Option{
		DSN: setting["dsn"].(string),
	}

	Add(name, option)
}

func AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		AddMap(name, setting.(map[string]interface{}))
	}
}

func Use(name string) *mongo.Client {
	if name == "" {
		name = "default"
	}

	if instance, ok := pool.Load(name); ok {
		return instance.(*mongo.Client)
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newClient, err := mongo.Connect(ctx, mongoOptions.Client().ApplyURI(option.DSN))
	if err != nil {
		panic("Failed to connect mongo " + name + " err: " + err.Error())
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = newClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic("Failed to connect mongo " + name + " err: " + err.Error())
	}

	if name == "default" {
		Default = newClient
	}

	pool.Store(name, newClient)
	return newClient
}
