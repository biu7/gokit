package log

import (
	"testing"
)

func TestHelper(t *testing.T) {
	l := NewSlog(LevelInfo, nil)
	l.Info("你好")

	l = With(l,
		"traceID", TraceID(),
		"caller", Caller(4),
	)

	l.Info("你好")
	Info("你好")
}
