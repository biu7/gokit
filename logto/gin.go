package logto

import (
	"errors"
	"github.com/biu7/gokit-qi/ginutils"
	"github.com/biu7/gokit-qi/ginutils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
)

func UserAuth(logto *Logto) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.AuthFail(c, ErrMissingToken)
			c.Abort()
			return
		}
		tokenStr = strings.TrimSuffix(tokenStr, "Bearer ")
		token, err := logto.Parse(tokenStr)
		if err != nil {
			response.AuthFail(c, ErrInvalidToken)
			c.Abort()
			return
		}
		if !token.Valid {
			response.AuthFail(c, ErrInvalidToken)
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		sub, err := claims.GetSubject()
		if err != nil {
			response.AuthFail(c, ErrInvalidToken)
			c.Abort()
			return
		}
		ginutils.SetStrUserID(c, sub)

		c.Next()
	}
}
