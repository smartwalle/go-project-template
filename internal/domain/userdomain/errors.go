package userdomain

import "project-template/internal/foundation/errors"

// 领域错误：通用语义（参数缺失、参数不合法等）直接用 foundation/errors 预置的错误，
// 需要区分业务语义时才在领域包里定义专属 code，code 带上业务前缀避免撞名。
var (
	// ErrUserNotFound 表示用户不存在。
	ErrUserNotFound = errors.New("user_not_found", "用户不存在")
)
