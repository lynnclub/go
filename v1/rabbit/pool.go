package rabbit

import (
	"fmt"
	"sync"

	"github.com/lynnclub/go/v1/logger"
	"github.com/lynnclub/go/v1/signal"
	"github.com/wagslane/go-rabbitmq"
)

var (
	pool           = &sync.Map{}             //实例池
	poolPublisher  = &sync.Map{}             //实例池
	poolConsumer   = &sync.Map{}             //实例池
	mutex          sync.Mutex                //互斥锁
	mutexPublisher sync.Mutex                //互斥锁
	mutexConsumer  sync.Mutex                //互斥锁
	options        = make(map[string]Option) //配置池
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

func Use(name string) *rabbitmq.Conn {
	if name == "" {
		name = "default"
	}

	if instance, ok := pool.Load(name); ok {
		return instance.(*rabbitmq.Conn)
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		if instance, ok = pool.Load(name); ok {
			return instance.(*rabbitmq.Conn)
		}
	}

	option, ok := options[name]
	if !ok {
		panic("Option not found " + name)
	}

	conn, err := rabbitmq.NewConn(option.DSN)
	if err != nil {
		panic("Failed to connect rabbitmq " + name + " err: " + err.Error())
	}

	pool.Store(name, conn)
	return conn
}

func GetPublisher(name string, optionFuncs ...func(*rabbitmq.PublisherOptions)) *rabbitmq.Publisher {
	if instance, ok := poolPublisher.Load(name); ok {
		return instance.(*rabbitmq.Publisher)
	} else {
		mutexPublisher.Lock()
		defer mutexPublisher.Unlock()
		if instance, ok = poolPublisher.Load(name); ok {
			return instance.(*rabbitmq.Publisher)
		}
	}

	publisher, err := rabbitmq.NewPublisher(Use(name), optionFuncs...)
	if err == nil {
		publisher.NotifyReturn(func(r rabbitmq.Return) {
			fmt.Println("rabbitmq message returned from server: " + string(r.Body))
			logger.Error("rabbitmq message returned from server: " + string(r.Body))
		})

		publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
			fmt.Println("rabbitmq publish", c.DeliveryTag)
		})
	} else {
		panic("Failed to new rabbitmq publisher " + name + " err: " + err.Error())
	}

	poolPublisher.Store(name, publisher)
	return publisher
}

func GetConsumer(name string, queue string, optionFuncs ...func(*rabbitmq.ConsumerOptions)) *rabbitmq.Consumer {
	if instance, ok := poolConsumer.Load(name); ok {
		return instance.(*rabbitmq.Consumer)
	} else {
		mutexConsumer.Lock()
		defer mutexConsumer.Unlock()
		if instance, ok = poolConsumer.Load(name); ok {
			return instance.(*rabbitmq.Consumer)
		}
	}

	consumer, err := rabbitmq.NewConsumer(Use(name), queue, optionFuncs...)
	if err != nil {
		panic("Failed to new rabbitmq publisher " + name + " err: " + err.Error())
	}

	poolConsumer.Store(name, consumer)
	return consumer
}

func AsyncCloseConsumer(consumer *rabbitmq.Consumer) {
	signal.Listen()

	go func() {
		fmt.Println("监听信号")
		<-signal.ChannelOS
		fmt.Println("收到信号", signal.Now)

		consumer.Close()
	}()
}
