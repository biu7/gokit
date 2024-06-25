package safe

import (
	"context"
	"github.com/biu7/gokit-qi/log"
	"runtime/debug"
)

func Go(ctx context.Context, fn func(ctx context.Context), logger log.Logger) {
	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("recover error", "error", r, "stack", string(debug.Stack()))
			}
		}()
		if fn != nil {
			fn(ctx)
		}
	}(ctx)
}
