package file

import (
	"encoding/csv"
	"os"

	"github.com/lynnclub/go/v1/pool"
)

func CSVWriter(filename string) *csv.Writer {
	file, err := os.Create(filename + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer
}

var CSVWriters = &pool.Pool[*csv.Writer]{
	Create: func(filename string) *csv.Writer {
		return CSVWriter(filename)
	},
}
