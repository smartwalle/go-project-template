package jsonrender

import (
	"github.com/smartwalle/rx"
)

// 统一的错误码。
const (
	// CodeInvalidJSON 表示请求体不是合法 JSON。
	CodeInvalidJSON = "invalid_json"

	// CodeRouteNotFound 表示路由不存在。
	CodeRouteNotFound = "route_not_found"
)

// Response 是统一响应体。
type Response struct {
	// Code 是结果码，失败时说明失败原因。
	Code string `json:"code"`

	// Message 是结果描述，失败时说明失败原因。
	Message string `json:"message"`

	// Data 是业务数据。
	Data any `json:"data,omitempty"`
}

// Success 写出成功响应。
func Success(c *rx.Context, status int, code, message string, data any) {
	c.JSON(status, Response{Code: code, Message: message, Data: data})
}

// Error 写出失败响应。
func Error(c *rx.Context, status int, code string, message string) {
	c.JSON(status, Response{Code: code, Message: message})
}
