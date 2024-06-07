package ginutils

import (
	"github.com/biu7/gokit-qi/log"
	"github.com/gin-gonic/gin"
	"strings"
)

type Middleware struct {
	logger log.Logger
}

func NewMiddleware(logger log.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func GetToken(c *gin.Context) string {
	tokenStr := c.Request.Header.Get("Authorization")
	if tokenStr == "" {
		return ""
	}
	return strings.TrimSuffix(tokenStr, "Bearer ")
}
