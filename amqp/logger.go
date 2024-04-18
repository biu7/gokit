package amqp

import (
	"fmt"
	"os"

	"github.com/biu7/gokit-qi/log"
	"github.com/wagslane/go-rabbitmq"
)

var _ rabbitmq.Logger = (*Logger)(nil)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) rabbitmq.Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Error("gorabbit: Fatal: " + fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error("gorabbit: Error: " + fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn("gorabbit: Warn: " + fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info("gorabbit: Info: " + fmt.Sprintf(format, v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug("gorabbit: Debug: " + fmt.Sprintf(format, v...))
}
