package http

import (
	"github.com/smartwalle/errors"
	_ "go-project-template/swagger"
)

var (
	ErrSuccess          = errors.New(100000, "成功")
	ErrInternalError    = errors.New(100001, "内部错误")
	ErrUnauthorized     = errors.New(100002, "未登录")
	ErrPermissionDenied = errors.New(100003, "没有操作权限")
	ErrInvalidSignature = errors.New(100004, "签名错误")
)

// Response 仅用作生成 Swagger 文档使用
type Response struct {
	Code    int32       `json:"code"`              // 错误码
	Message string      `json:"message,omitempty"` // 错误消息
	Data    interface{} `json:"data,omitempty"`    // 数据
}
