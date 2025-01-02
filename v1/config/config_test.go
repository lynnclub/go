package config

import (
	"os"
	"path"
	"runtime"
	"testing"
)

// TestStart 启动
func TestStart(t *testing.T) {
	// 不在工作目录下运行，获取不到正确目录
	// pathDir, _ := os.Getwd()

	// 通用方式
	_, file, _, _ := runtime.Caller(0)
	pathDir := path.Dir(file)
	Start("_TEST_MODE", pathDir+"/config")
	env := Viper.GetString("name")
	if env != "release" {
		panic("config read error")
	}

	err := os.Setenv("_TEST_MODE", "test")
	if err != nil {
		panic(err.Error())
	}

	Viper = nil
	Start("_TEST_MODE", pathDir+"/config")
	env = Viper.GetString("name")
	if env != "test" {
		panic("config read error")
	}

	Viper = nil
	Start("_TEST_MODE", "not")
	env = Viper.GetString("name")
	if env != "test" {
		panic("config read error")
	}
}
