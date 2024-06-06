package ginutils

import "github.com/gin-gonic/gin"

const (
	ContextUserID   = "user_id"
	ContextDeviceID = "device_id"
	ContextPlatform = "platform"
	ContextResponse = "gin_response"
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
