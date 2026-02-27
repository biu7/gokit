package amqp

import (
	"fmt"
	"os"

	"github.com/biu7/gokit/log"
	"github.com/wagslane/go-rabbitmq"
)

var _ rabbitmq.Logger = (*logAdapter)(nil)

type logAdapter struct {
	logger log.Logger
}

func newLogAdapter(logger log.Logger) rabbitmq.Logger {
	return &logAdapter{logger: logger}
}

func (l *logAdapter) Fatalf(format string, v ...interface{}) {
	l.logger.Error("gorabbit: Fatal: " + fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logAdapter) Errorf(format string, v ...interface{}) {
	l.logger.Error("gorabbit: Error: " + fmt.Sprintf(format, v...))
}

func (l *logAdapter) Warnf(format string, v ...interface{}) {
	l.logger.Warn("gorabbit: Warn: " + fmt.Sprintf(format, v...))
}

func (l *logAdapter) Infof(format string, v ...interface{}) {
	l.logger.Info("gorabbit: Info: " + fmt.Sprintf(format, v...))
}

func (l *logAdapter) Debugf(format string, v ...interface{}) {
	l.logger.Debug("gorabbit: Debug: " + fmt.Sprintf(format, v...))
}
