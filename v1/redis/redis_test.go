package redis

import (
	"os"
	"testing"

	"github.com/lynnclub/go/v1/config"
)

// TestGoRedis go-redis
func TestGoRedis(t *testing.T) {
	err := os.Setenv("_TEST_MODE", "test")
	if err != nil {
		panic(err.Error())
	}

	config.Start("_TEST_MODE", "../config")
	AddMapBatch(config.Viper.GetStringMap("redis"))

	db := Use("")
	if err = db.Ping(Ctx).Err(); err != nil {
		panic("redis go-redis error " + err.Error())
	}
	if err = Default.Ping(Ctx).Err(); err != nil {
		panic("redis go-redis error " + err.Error())
	}

	_, err = Default.Get(Ctx, "the_key_does_not_exist_yeah").Result()
	if err != Nil {
		panic("redis go-redis error " + err.Error())
	}
}
