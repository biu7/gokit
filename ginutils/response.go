package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"net/http"
)

const (
	ctxKeyResponse = "gin_response"
)

var marshaller = protojson.MarshalOptions{
	AllowPartial:      true,
	UseEnumNumbers:    true,
	EmitDefaultValues: true,
}

func ProtoJSON(c *gin.Context, code int, data proto.Message, msg string) {
	// for logging
	c.Set(ctxKeyResponse, &CommonResponse{
		Code:    int32(code),
		Message: msg,
	})

	var resp = &CommonResponse{
		Code:    int32(code),
		Message: msg,
	}
	if data != nil {
		resp.Data, _ = anypb.New(data)
	}
	b, _ := marshaller.Marshal(resp)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Render(http.StatusOK, render.String{
		Format: "%s",
		Data:   []any{string(b)},
	})
}
