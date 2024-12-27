package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	CodeOK        = http.StatusOK
	CodeError     = http.StatusBadRequest
	CodeAuthError = http.StatusUnauthorized
)

const (
	ContextResponse = "gin_response"
)

var marshaller = protojson.MarshalOptions{
	AllowPartial:      true,
	UseEnumNumbers:    true,
	EmitDefaultValues: true,
}

func ProtoJSON(c *gin.Context, code int, data proto.Message, msg string) {
	resp := &CommonResponse{
		Code:    int32(code),
		Message: msg,
	}
	if data != nil {
		anyData, _ := anypb.New(data)
		resp.Data = anyData
	}

	// for logging
	SetResponseStatus(c, resp.GetCode(), resp.GetMessage())

	b, _ := marshaller.Marshal(resp)
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
	return commonResp.GetCode(), commonResp.GetMessage()
}

func Success(c *gin.Context, body proto.Message) {
	ProtoJSON(c, CodeOK, body, "")
}

func Fail(c *gin.Context, err error) {
	_, message := errMessage(err)
	ProtoJSON(c, CodeError, nil, message)
}

func AuthFail(c *gin.Context, err error) {
	_, message := errMessage(err)

	ProtoJSON(c, CodeAuthError, nil, message)
}
