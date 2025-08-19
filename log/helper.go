package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// Valuer is returns a log value.
type Valuer func(ctx context.Context) interface{}

// Value return the function value.
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// TraceID returns a traceid valuer.
func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return span.TraceID().String()
		}

		if c, ok := ctx.(*gin.Context); ok {
			if span := trace.SpanContextFromContext(c.Request.Context()); span.HasTraceID() {
				return span.TraceID().String()
			}
		}
		return ""
	}
}

// SpanID returns a spanid valuer.
func SpanID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return span.SpanID().String()
		}

		if c, ok := ctx.(*gin.Context); ok {
			if span := trace.SpanContextFromContext(c.Request.Context()); span.HasSpanID() {
				return span.SpanID().String()
			}
		}
		return ""
	}
}

func bindValues(ctx context.Context, args []any) {
	for i := 1; i < len(args); i += 2 {
		if v, ok := args[i].(Valuer); ok {
			args[i] = v(ctx)
		}
	}
}

func containsValuer(args []any) bool {
	for i := 1; i < len(args); i += 2 {
		if _, ok := args[i].(Valuer); ok {
			return true
		}
	}
	return false
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}
