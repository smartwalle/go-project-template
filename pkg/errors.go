package pkg

import "github.com/smartwalle/errors"

// 通用错误
var (
	ErrSuccess          = errors.New(100000, "成功")
	ErrInternalError    = errors.New(100001, "内部错误")
	ErrUnauthorized     = errors.New(100002, "未登录")
	ErrPermissionDenied = errors.New(100003, "没有操作权限")
)
