package db

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/lynnclub/go/v1/config"
)

// TestGORM GORM
func TestGORM(t *testing.T) {
	err := os.Setenv("_TEST_MODE", "test")
	if err != nil {
		panic(err.Error())
	}

	config.Start("_TEST_MODE", "../config/config")

	AddMapBatch(config.Viper.GetStringMap("db"))

	var wg sync.WaitGroup
	for loop := 0; loop < 10; loop++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			goDB := Use("")
			fmt.Printf("%p\n", goDB)
			if goDB != Default {
				panic("GORM mysql not reuse")
			}
		}()
	}
	wg.Wait()

	db, _ := Use("").DB()
	if err = db.Ping(); err != nil {
		panic("GORM mysql error " + err.Error())
	}
	db, _ = Default.DB()
	if err = db.Ping(); err != nil {
		panic("GORM mysql error " + err.Error())
	}

	db, _ = Use("postgres").DB()
	if err = db.Ping(); err != nil {
		panic("GORM postgres error " + err.Error())
	}
}
