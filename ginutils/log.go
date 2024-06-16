package ginutils

import (
	"bytes"
	"github.com/biu7/gokit-qi/ginutils/response"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ValuerFunc func(c *gin.Context) (string, any)

func (g *Middleware) Log(values ...ValuerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		buf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		var (
			status = int32(c.Writer.Status())
			msg    string
		)

		respStatus, message := response.GetResponseStatus(c)
		if respStatus != 0 {
			status = respStatus
			msg = message
		}

		fields := []any{
			"time", end.Format(time.RFC3339),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"message", msg,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"latency", latency,
		}
		var headers = make(map[string][]string)
		for k, v := range c.Request.Header {
			if strings.HasPrefix(strings.ToLower(k), "x-wx") {
				headers[k] = v
			}
		}
		// 添加请求头日志
		if len(headers) > 0 {
			fields = append(fields, "header", headers)
		}
		// 添加自定义日志
		for _, v := range values {
			key, val := v(c)
			fields = append(fields, key, val)
		}
		// 添加 body 日志
		if len(buf) < 1024 {
			fields = append(fields, "body", string(buf))
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				errField := append(fields, "error", e)
				g.logger.Ctx(c).Error("RecordAPI with error", errField...)
			}
		} else {
			g.logger.Ctx(c).Info("RecordAPI", fields...)
		}
	}
}

func (g *Middleware) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, X-Meokii-Openid")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		c.Next()
	}
}
