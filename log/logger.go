package log

import (
	"context"
	"log/slog"
)

type Level slog.Level

const (
	LevelDebug Level = Level(slog.LevelDebug)
	LevelInfo  Level = Level(slog.LevelInfo)
	LevelWarn  Level = Level(slog.LevelWarn)
	LevelError Level = Level(slog.LevelError)
)

type Logger interface {
	Ctx(ctx context.Context) Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func NewLogger(level Level) Logger {
	return NewSlog(level, nil)
}

type logger struct {
	logger    Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (c *logger) Ctx(ctx context.Context) Logger {
	return WithContext(ctx, c)
}

func (c *logger) Debug(msg string, args ...any) {
	c.logger.Debug(msg, c.bindArgs(args)...)
}

func (c *logger) Info(msg string, args ...any) {
	c.logger.Info(msg, c.bindArgs(args)...)
}

func (c *logger) Warn(msg string, args ...any) {
	c.logger.Warn(msg, c.bindArgs(args)...)
}

func (c *logger) Error(msg string, args ...any) {
	c.logger.Error(msg, c.bindArgs(args)...)
}

func (c *logger) bindArgs(args []any) []any {
	kvs := make([]any, 0, len(c.prefix)+len(args))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(c.ctx, kvs)
	}
	kvs = append(kvs, args...)
	return kvs
}

func With(l Logger, args ...any) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: args, hasValuer: containsValuer(args), ctx: context.Background()}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(args))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, args...)
	return &logger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

func WithContext(ctx context.Context, l Logger) Logger {
	c, ok := l.(*logger)
	if ok {
		return &logger{
			logger:    c.logger,
			prefix:    c.prefix,
			hasValuer: c.hasValuer,
			ctx:       ctx,
		}
	}
	return &logger{logger: l, ctx: ctx}
}
