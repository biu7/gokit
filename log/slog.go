package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var _ Logger = (*slogLogger)(nil)

type slogLogger struct {
	logger *slog.Logger
}

func newSlogLogger(level Level, writer io.Writer) Logger {
	if writer == nil {
		writer = os.Stderr
	}
	return &slogLogger{
		logger: slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: slog.Level(level),
		})),
	}
}

func (s *slogLogger) log(level slog.Level, msg string, args ...any) {
	if !s.logger.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = s.logger.Handler().Handle(context.Background(), r)
}

func (s *slogLogger) Debug(msg string, args ...any) {
	s.log(slog.LevelDebug, msg, args...)
}

func (s *slogLogger) Info(msg string, args ...any) {
	s.log(slog.LevelInfo, msg, args...)
}

func (s *slogLogger) Warn(msg string, args ...any) {
	s.log(slog.LevelWarn, msg, args...)
}

func (s *slogLogger) Error(msg string, args ...any) {
	s.log(slog.LevelError, msg, args...)
}

func (s *slogLogger) Ctx(ctx context.Context) Logger {
	return WithContext(ctx, s)
}
