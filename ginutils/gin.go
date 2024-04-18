package ginutils

import "github.com/biu7/gokit-qi/log"

type Middleware struct {
	logger log.Logger
}

func NewMiddleware(logger log.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}
