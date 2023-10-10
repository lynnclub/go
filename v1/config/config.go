package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Env   = "dev"      //环境 建议dev开发、qa测试、yc压测、prv预览、pro生产
	Viper *viper.Viper //配置
)

// Start 启动
func Start(envKey, path string) {
	if Viper != nil {
		return
	}

	Env = os.Getenv(envKey)
	if Env == "" {
		if flagMode := flag.Lookup("m"); flagMode == nil {
			input := flag.String("m", "dev", "环境")
			flag.Parse()
			Env = *input
		} else {
			Env = flagMode.Value.String()
		}
	}

	file := path + "/" + Env + ".yaml"

	Viper = viper.New()
	Viper.SetConfigFile(file)
	if err := Viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	Viper.WatchConfig()

	fmt.Println("config start:", file)
}
