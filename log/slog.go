package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var _ Logger = (*Slog)(nil)

var slogLevels = map[Level]slog.Level{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
}

type Slog struct {
	logger    *slog.Logger
	addSource bool
}

func NewSlog(level Level, writer io.Writer) Logger {
	if writer == nil {
		writer = os.Stderr
	}
	return &Slog{
		logger: slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
			//AddSource: true,
			Level: slog.Level(level),
		})),
		addSource: false,
	}
}

func (s *Slog) level(level Level) slog.Level {
	lv, ok := slogLevels[level]
	if ok {
		return lv
	}
	return slog.LevelInfo
}

func (s *Slog) log(level slog.Level, msg string, args ...any) {
	if !s.logger.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = s.logger.Handler().Handle(context.Background(), r)
}

func (s *Slog) Debug(msg string, args ...any) {
	s.log(slog.LevelDebug, msg, args...)
}

func (s *Slog) Info(msg string, args ...any) {
	s.log(slog.LevelInfo, msg, args...)
}

func (s *Slog) Warn(msg string, args ...any) {
	s.log(slog.LevelWarn, msg, args...)
}

func (s *Slog) Error(msg string, args ...any) {
	s.log(slog.LevelError, msg, args...)
}

func (s *Slog) Ctx(ctx context.Context) Logger {
	return WithContext(ctx, s)
}
