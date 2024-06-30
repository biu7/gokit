package response

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"net/http"
)

const (
	ContextResponse = "gin_response"
)

const (
	CodeOK        = http.StatusOK
	CodeError     = http.StatusBadRequest
	CodeAuthError = http.StatusUnauthorized
)

var marshaller = protojson.MarshalOptions{
	AllowPartial:      true,
	UseEnumNumbers:    true,
	EmitDefaultValues: true,
}

func ProtoJson(c *gin.Context, code int, data proto.Message, msg string) {
	// 日志记录
	SetResponseStatus(c, code, msg)

	var anyData *anypb.Any
	if data != nil {
		anyData, _ = anypb.New(data)
	}
	resp := &CommonResponse{
		Code:    int32(code),
		Message: msg,
		Data:    anyData,
	}
	b, _ := marshaller.Marshal(resp)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Render(http.StatusOK, render.String{
		Format: "%s",
		Data:   []any{string(b)},
	})
}

func Success(c *gin.Context, body proto.Message) {
	ProtoJson(c, CodeOK, body, "success")
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

func Fail(c *gin.Context, err error) {
	code, message := errMessage(err)
	ProtoJson(c, code, nil, message)
}

func SetResponseStatus(c *gin.Context, code int, msg string) {
	c.Set(ContextResponse, &CommonResponse{
		Code:    int32(code),
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
	return commonResp.GetCode(), commonResp.GetMessage()
}
