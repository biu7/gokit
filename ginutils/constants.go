package ginutils

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const (
	ContextUserID   = "user_id"
	ContextDeviceID = "device_id"
	ContextPlatform = "platform"
	ContextLang     = "lang"
)

func SetStrUserID(c *gin.Context, userID string) {
	c.Set(ContextUserID, userID)
}

func SetStrDeviceID(c *gin.Context, deviceID string) {
	c.Set(ContextDeviceID, deviceID)
}

func SetStrPlatform(c *gin.Context, platform string) {
	c.Set(ContextPlatform, platform)
}

func GetStrUserID(c *gin.Context) string {
	return c.GetString(ContextUserID)
}

func GetStrDeviceID(c *gin.Context) string {
	return c.GetString(ContextDeviceID)
}

func GetStrPlatform(c *gin.Context) string {
	return c.GetString(ContextPlatform)
}

func SetIntUserID(c *gin.Context, userID int64) {
	c.Set(ContextUserID, userID)
}

func GetIntUserID(c *gin.Context) int64 {
	return c.GetInt64(ContextUserID)
}

func SetLanguage(c *gin.Context, lang language.Tag) {
	c.Set(ContextLang, lang)
}

func GetContextLang(ctx context.Context, def language.Tag) language.Tag {
	if v := ctx.Value(ContextLang); v != nil {
		if lang, ok := v.(language.Tag); ok {
			return lang
		}
	}
	return def
}
