package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	code     int
	message  string
	httpCode int
}

func New(httpCode, bizCode int, message string) *Error {
	return &Error{
		code:     bizCode,
		message:  message,
		httpCode: httpCode,
	}
}

func Newf(httpCode, bizCode int, format string, args ...any) *Error {
	return &Error{
		code:     bizCode,
		message:  fmt.Sprintf(format, args...),
		httpCode: httpCode,
	}
}

func BadRequest(bizCode int, message string) *Error {
	return New(http.StatusBadRequest, bizCode, message)
}

func NotFound(bizCode int, message string) *Error {
	return New(http.StatusNotFound, bizCode, message)
}

func Unauthorized(bizCode int, message string) *Error {
	return New(http.StatusUnauthorized, bizCode, message)
}

func Internal(bizCode int, message string) *Error {
	return New(http.StatusInternalServerError, bizCode, message)
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) HTTPCode() int {
	return e.httpCode
}

func (e *Error) BizCode() int {
	return e.code
}

func (e *Error) WithMessage(message string) *Error {
	return &Error{
		code:     e.code,
		message:  message,
		httpCode: e.httpCode,
	}
}

func (e *Error) WithMessagef(format string, args ...any) *Error {
	return &Error{
		code:     e.code,
		message:  fmt.Sprintf(format, args...),
		httpCode: e.httpCode,
	}
}
