package redis

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/lynnclub/go/v1/config"
)

// TestGoRedis go-redis
func TestGoRedis(t *testing.T) {
	err := os.Setenv("_TEST_MODE", "test")
	if err != nil {
		panic(err.Error())
	}

	config.Start("_TEST_MODE", "../config/config")
	AddMapBatch(config.Viper.GetStringMap("redis"))

	var wg sync.WaitGroup
	for loop := 0; loop < 10; loop++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			goDB := Use("")
			fmt.Printf("%p\n", goDB)
			if goDB != Use("default") {
				panic("redis not reuse")
			}
		}()
	}
	wg.Wait()

	db := Use("")
	if err = db.Ping(Ctx).Err(); err != nil {
		panic("redis go-redis error " + err.Error())
	}

	_, err = db.Get(Ctx, "the_key_does_not_exist_yeah").Result()
	if err != Nil {
		panic("redis go-redis error " + err.Error())
	}
}
