package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-template/internal/tool"
)

type Request interface {
	Check() error
}

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	kDefaultErrorCode    = "operation_failed"
	kDefaultErrorMessage = "操作失败"
)

func NewResponse(code, message string, data interface{}) *Response {
	var r = &Response{}
	r.Code = code
	r.Message = message
	r.Data = data
	return r
}

func Handle[T Request, R interface{}](handler func(c *gin.Context, req T) (R, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		if handler != nil {
			var resp, err = handle(c, handler)
			if err != nil {
				var code = kDefaultErrorCode
				var message = kDefaultErrorMessage

				var rErr tool.Error
				if errors.As(err, &rErr) {
					code = rErr.Code()
					message = rErr.Message()
				}
				c.JSON(http.StatusBadRequest, NewResponse(code, message, nil))
				return
			}
			c.JSON(http.StatusOK, NewResponse("", "", resp))
		}
	}
}

func handle[T Request, R interface{}](c *gin.Context, handler func(c *gin.Context, req T) (R, error)) (resp R, err error) {
	var req T
	if err = c.BindJSON(&req); err != nil {
		return resp, err
	}

	if err = req.Check(); err != nil {
		return resp, err
	}

	if resp, err = handler(c, req); err != nil {
		return resp, err
	}
	return resp, nil
}
