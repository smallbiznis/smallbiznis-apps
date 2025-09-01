package errutil

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BaseError struct {
	Code    CoreStatus `json:"code"`
	Message string     `json:"message"`
	Details []Detail   `json:"details,omitempty"`
	Err     error      `json:"-"`
}

func (e BaseError) URL() string {
	values := url.Values{}

	values.Set("error_code", string(e.Code))
	values.Set("error_message", e.Message)

	for _, d := range e.Details {
		values.Set("details["+strings.TrimSpace(d.Field)+"]", d.Message)
	}

	return values.Encode()
}

func (e BaseError) JSON() gin.H {
	return gin.H{
		"error": gin.H{
			"code":    e.Message,
			"message": e.Err.Error(),
			"details": e.Details,
		},
	}
}

func (e BaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code CoreStatus, message string, err error, details ...Detail) error {
	return BaseError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

func NotFound(msg string, err error, details ...Detail) error {
	return New(StatusNotFound, msg, err, details...)
}

func UnprocessableEntity(msg string, err error, details ...Detail) error {
	return New(StatusUnprocessableEntity, msg, err, details...)
}

func UnsupportedMediaType(msg string, err error, details ...Detail) error {
	return New(StatusUnsupportedMediaType, msg, err, details...)
}

func Conflict(msg string, err error, details ...Detail) error {
	return New(StatusConflict, msg, err, details...)
}

func BadRequest(msg string, err error, details ...Detail) error {
	return New(StatusBadRequest, msg, err, details...)
}

func ValidationFailed(msg string, err error, details ...Detail) error {
	return New(StatusValidationFailed, msg, err, details...)
}

func Internal(msg string, err error, details ...Detail) error {
	return New(StatusInternal, msg, err, details...)
}

func Timeout(msg string, err error, details ...Detail) error {
	return New(StatusTimeout, msg, err, details...)
}

func Unauthorized(msg string, err error, details ...Detail) error {
	return New(StatusUnauthorized, msg, err, details...)
}

func Forbidden(msg string, err error, details ...Detail) error {
	return New(StatusForbidden, msg, err, details...)
}

func TooManyRequest(msg string, err error, details ...Detail) error {
	return New(StatusTooManyRequests, msg, err, details...)
}

func ClientClosedRequest(msg string, err error, details ...Detail) error {
	return New(StatusClientClosedRequest, msg, err, details...)
}

func NotImplemented(msg string, err error, details ...Detail) error {
	return New(StatusNotImplemented, msg, err, details...)
}

func BadGateway(msg string, err error, details ...Detail) error {
	return New(StatusBadGateway, msg, err, details...)
}
