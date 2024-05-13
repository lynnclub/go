package db

import (
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

	config.Start("_TEST_MODE", "../config")

	AddMapBatch(config.Viper.GetStringMap("db"))

	db, _ := Use("").DB()
	if err = db.Ping(); err != nil {
		panic("GORM mysql error " + err.Error())
	}
	db, _ = Default.DB()
	if err = db.Ping(); err != nil {
		panic("GORM mysql error " + err.Error())
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			goDb := Use("")
			goDb2 := Use("")
			if *goDb != *goDb2 {
				panic("GORM mysql not reuse")
			}
		}()
	}
	wg.Wait()

	db, _ = Use("postgres").DB()
	if err = db.Ping(); err != nil {
		panic("GORM postgres error " + err.Error())
	}
}
