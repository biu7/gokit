package response

import (
	"net/http"

	"github.com/biu7/gokit/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/proto"
)

const (
	CodeOK        = 0
	CodeError     = http.StatusBadRequest
	CodeAuthError = http.StatusUnauthorized
)

const (
	ContextResponse = "gin_response"
)

type CommonResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ProtoJSON(c *gin.Context, code int, data proto.Message, msg string) {
	resp := &CommonResponse{
		Code:    int32(code),
		Message: msg,
	}
	if data != nil {
		resp.Data = data
	}

	// for logging
	SetResponseStatus(c, resp.Code, resp.Message)

	b, _ := json.Marshal(resp)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Render(http.StatusOK, render.String{
		Format: "%s",
		Data:   []any{string(b)},
	})
}

func errMessage(err error) (int, string) {
	if kErr := errors.FromError(err); kErr != nil {
		if kErr.GetCode() != errors.UnknownCode && kErr.GetCode() != 0 {
			return int(kErr.GetCode()), kErr.GetMessage()
		}
		return CodeError, kErr.GetMessage()
	}
	return CodeError, "unknown error"
}

func SetResponseStatus(c *gin.Context, code int32, msg string) {
	c.Set(ContextResponse, &CommonResponse{
		Code:    code,
		Message: msg,
	})
}

func GetResponseStatus(c *gin.Context) (int32, string) {
	respStatus, ok := c.Get(ContextResponse)
	if !ok {
		return 0, ""
	}
	commonResp, ok := respStatus.(*CommonResponse)
	if !ok {
		return 0, ""
	}
	return commonResp.Code, commonResp.Message
}

func Success(c *gin.Context, body proto.Message) {
	ProtoJSON(c, CodeOK, body, "ok")
}

func Fail(c *gin.Context, err error) {
	_, message := errMessage(err)
	ProtoJSON(c, CodeError, nil, message)
}

func AuthFail(c *gin.Context, err error) {
	_, message := errMessage(err)

	ProtoJSON(c, CodeAuthError, nil, message)
}
