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
	return newSlogLogger(level, nil)
}

type logger struct {
	inner     Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (l *logger) Ctx(ctx context.Context) Logger {
	return WithContext(ctx, l)
}

func (l *logger) Debug(msg string, args ...any) {
	l.inner.Debug(msg, l.bindArgs(args)...)
}

func (l *logger) Info(msg string, args ...any) {
	l.inner.Info(msg, l.bindArgs(args)...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.inner.Warn(msg, l.bindArgs(args)...)
}

func (l *logger) Error(msg string, args ...any) {
	l.inner.Error(msg, l.bindArgs(args)...)
}

func (l *logger) bindArgs(args []any) []any {
	kvs := make([]any, 0, len(l.prefix)+len(args))
	kvs = append(kvs, l.prefix...)
	if l.hasValuer {
		bindValues(l.ctx, kvs)
	}
	kvs = append(kvs, args...)
	return kvs
}

func With(l Logger, args ...any) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{inner: l, prefix: args, hasValuer: containsValuer(args), ctx: context.Background()}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(args))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, args...)
	return &logger{
		inner:     c.inner,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

func WithContext(ctx context.Context, l Logger) Logger {
	c, ok := l.(*logger)
	if ok {
		return &logger{
			inner:     c.inner,
			prefix:    c.prefix,
			hasValuer: c.hasValuer,
			ctx:       ctx,
		}
	}
	return &logger{inner: l, ctx: ctx}
}
