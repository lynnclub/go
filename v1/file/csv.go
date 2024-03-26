package file

import (
	"encoding/csv"
	"os"
)

func CSVWriter(filename string, headers ...string) *csv.Writer {
	file, err := os.Create(filename + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if len(headers) > 0 {
		writer.Write(headers)
	}

	return writer
}
