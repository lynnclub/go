package config

import (
	"os"
	"testing"
)

// TestStart 启动
func TestStart(t *testing.T) {
	path, _ := os.Getwd()
	Start("_TEST_MODE", path)
	env := Viper.GetString("name")
	if env != "release" {
		panic("config read error")
	}

	Viper = nil
	err := os.Setenv("_TEST_MODE", "test")
	if err != nil {
		panic(err.Error())
	}

	Start("_TEST_MODE", path)
	env = Viper.GetString("name")
	if env != "test" {
		panic("config read error")
	}
}
