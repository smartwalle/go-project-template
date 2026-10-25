package userdomain

import "project-template/internal/platform/errors"

// 领域错误：通用语义（资源不存在、参数缺失、参数不合法等）直接用 platform/errors 预置的错误码，
// 业务专属的错误码才定义在领域包里。
var (
	// ErrNotFound 表示目标用户不存在。
	ErrNotFound = errors.NotFound("user not found")
)
