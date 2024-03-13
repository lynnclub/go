package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	Env      = "release"  //环境 建议dev开发、test测试、release生产
	Viper    *viper.Viper //配置
	BasePath string       //根路径
)

// Start 启动
func Start(envKey, path string) {
	if Viper != nil {
		return
	}

	var env string
	if flagMode := flag.Lookup("e"); flagMode == nil {
		input := flag.String("e", "", "环境")
		flag.Parse()
		env = *input
	} else {
		env = flagMode.Value.String()
	}

	if env == "" {
		env = os.Getenv(envKey)
	}

	if env != "" {
		Env = env
	}

	file := path + "/" + Env + ".yaml"
	BasePath = filepath.Dir(path)

	Viper = viper.New()
	Viper.SetConfigFile(file)
	if err := Viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	fmt.Println("config start:", file)
}
