package redis

import (
	"context"
	"log"
	"os"
)

// stdoutLogger go-redis 默认输出到 stderr，纠正为 stdout
type stdoutLogger struct {
	log *log.Logger
}

func (l *stdoutLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

// newStdoutLogger 输出到 stdout
func newStdoutLogger() *stdoutLogger {
	return &stdoutLogger{
		log: log.New(os.Stdout, "redis: ", log.LstdFlags|log.Lshortfile),
	}
}
