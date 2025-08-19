package ctxvalue

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cast"
	"golang.org/x/text/language"
)

const (
	ContextUserID       = "user_id"
	ContextDeviceID     = "device_id"
	ContextPlatform     = "platform"
	ContextLang         = "lang"
	ContextWechatOpenID = "wechat_openid"
	ContextVersion      = "version"
)

func SetUserID(ctx context.Context, userID int64) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(ContextUserID, userID)
		return ginCtx
	}
	return context.WithValue(ctx, ContextUserID, userID)
}

func GetUserID(ctx context.Context) int64 {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		val, exists := ginCtx.Get(ContextUserID)
		if exists {
			return cast.ToInt64(val)
		}
	}
	return cast.ToInt64(ctx.Value(ContextUserID))
}

func SetDeviceID(ctx context.Context, deviceID string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(ContextDeviceID, deviceID)
		return ginCtx
	}
	return context.WithValue(ctx, ContextDeviceID, deviceID)
}

func GetDeviceID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		val, exists := ginCtx.Get(ContextDeviceID)
		if exists {
			return cast.ToString(val)
		}
	}
	return cast.ToString(ctx.Value(ContextDeviceID))
}

func SetPlatform(ctx context.Context, platform string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(ContextPlatform, platform)
		return ginCtx
	}
	return context.WithValue(ctx, ContextPlatform, platform)
}

func GetContextWechatOpenID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		val, exists := ginCtx.Get(ContextWechatOpenID)
		if exists {
			return cast.ToString(val)
		}
	}
	return cast.ToString(ctx.Value(ContextWechatOpenID))
}

func SetContextWechatOpenID(ctx context.Context, openid string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(ContextWechatOpenID, openid)
		return ginCtx
	}
	return context.WithValue(ctx, ContextWechatOpenID, openid)
}

func GetPlatform(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		val, exists := ginCtx.Get(ContextPlatform)
		if exists {
			return cast.ToString(val)
		}
	}
	return cast.ToString(ctx.Value(ContextPlatform))
}

func SetLanguage(ctx context.Context, lang language.Tag) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(ContextLang, lang)
		return ginCtx
	}
	return context.WithValue(ctx, ContextLang, lang)
}

func GetContextLang(ctx context.Context, def language.Tag) language.Tag {
	var val any
	if ginCtx, ok := ctx.(*gin.Context); ok {
		v, exists := ginCtx.Get(ContextLang)
		if exists {
			val = v
		}
	}

	if val == nil {
		val = ctx.Value(ContextLang)
	}

	if lang, ok := val.(language.Tag); ok {
		return lang
	}
	return def
}

func SetContextVersion(ctx context.Context, version int32) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(ContextVersion, version)
		return ginCtx
	}
	return context.WithValue(ctx, ContextVersion, version)
}

func GetContextVersion(ctx context.Context) int32 {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		val, exists := ginCtx.Get(ContextVersion)
		if exists {
			return cast.ToInt32(val)
		}
	}
	return cast.ToInt32(ctx.Value(ContextVersion))
}
