package userdomain

import "project-template/internal/foundation/errors"

var (
	// ErrNotFound 表示用户不存在。
	// 使用通用的 not_found 错误码，错误信息包含用户领域语义。
	ErrNotFound = errors.NotFound("user not found")
)
