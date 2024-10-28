package file

import (
	"encoding/csv"
	"os"
)

func CSVWriter(filename string, overwrite bool, headers ...string) *csv.Writer {
	var file *os.File
	var err error

	if overwrite {
		file, err = os.Create(filename + ".csv") // 覆盖模式
	} else {
		file, err = os.OpenFile(filename+".csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}

	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(file)

	if overwrite || isEmpty(file) {
		writer.Write(headers)
	}

	return writer
}

func isEmpty(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	return stat.Size() == 0
}
