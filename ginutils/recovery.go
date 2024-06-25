package ginutils

import (
	"errors"
	"github.com/biu7/gokit-qi/ginutils/response"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func netOPError(err error) bool {
	var opError *net.OpError
	return errors.As(err, &opError)
}

func sysCallError(err error) bool {
	var se *os.SyscallError
	return errors.As(err, &se)
}

func (g *Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				if err, ok := rec.(error); ok {
					var brokenPipe bool
					if netOPError(err) && sysCallError(err) {
						if strings.Contains(strings.ToLower(err.Error()), "broken pipe") || strings.Contains(strings.ToLower(err.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
					if brokenPipe {
						g.logger.Ctx(c).Error("broken pipe error",
							"path", c.Request.URL.Path,
							"error", err,
							"request", string(httpRequest))
						// If the connection is dead, we can't write a statusfrtgdk,,,,,,,,,m to it.
						_ = c.Error(err.(error))
						c.Abort()
						return
					}
				}

				g.logger.Ctx(c).Error("[Recovery from panic]",
					"time", time.Now(),
					"error", rec,
					"request", string(httpRequest))
				defaultHandleRecovery(c, rec)
			}
		}()
		c.Next()
	}
}

func defaultHandleRecovery(c *gin.Context, err interface{}) {
	if _, ok := err.(error); ok {
		response.Fail(c, err.(error))
		c.Abort()
		return
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
