package ctxvalue

import (
	"context"
	"github.com/spf13/cast"
	"golang.org/x/text/language"
)

const (
	ContextUserID   = "user_id"
	ContextDeviceID = "device_id"
	ContextPlatform = "platform"
	ContextLang     = "lang"
)

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, ContextUserID, userID)
}

func GetUserID(ctx context.Context) int64 {
	return cast.ToInt64(ctx.Value(ContextUserID))
}

func SetDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, ContextDeviceID, deviceID)
}

func GetDeviceID(ctx context.Context) string {
	return cast.ToString(ctx.Value(ContextDeviceID))
}

func SetPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, ContextPlatform, platform)
}

func GetPlatform(ctx context.Context) string {
	return cast.ToString(ctx.Value(ContextPlatform))
}

func SetLanguage(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, ContextLang, lang)
}

func GetContextLang(ctx context.Context, def language.Tag) language.Tag {
	if v := ctx.Value(ContextLang); v != nil {
		if lang, ok := v.(language.Tag); ok {
			return lang
		}
	}
	return def
}
