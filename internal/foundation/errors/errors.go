// Package errors 提供通用的错误类型与常用错误码。
//
// Error 用 code 表达语义、用 message 描述细节，调用方按 code 判断错误类型。
//
// 本包只放跨业务通用的语义，不放业务专属的错误。
// 业务代码直接复用这里的构造函数，message 写清楚是哪个资源：
//
//	errors.NotFound("用户不存在")
//	errors.NotFound("订单不存在")
//
// 同一语义用同一个 code，调用方按 code 判断，message 只描述细节。
// 需要区分业务语义时，在各自的领域包里定义专属的 code，用 New 构造，
// code 带上业务前缀避免撞名：
//
//	errors.New("user_not_found", "用户不存在")
//	errors.New("order_not_found", "订单不存在")
//
// 两种写法都遵守同样的约定：code 稳定可读、message 描述细节。
package errors

import (
	stderrors "errors"
	"strings"
)

// Error 是业务错误。
//
// 它是值类型，code 与 message 相同的两个 Error 相等，
// 可以直接用 errors.Is 比较。
//
// message 允许随业务变化（"用户不存在"、"订单不存在"），
// 因此判断错误语义要按 Code 比，不要按值比。
type Error struct {
	code    string
	message string
}

// New 构造一个业务错误。
func New(code string, message string) Error {
	return Error{code: code, message: message}
}

func (err Error) Error() string {
	return strings.Join([]string{err.code, err.message}, " - ")
}

func (err Error) Code() string {
	return err.code
}

func (err Error) Message() string {
	return err.message
}

// As 判断错误链里有没有目标错误，用法同标准库 errors.As。
func As(err error, target any) bool {
	return stderrors.As(err, target)
}

// Is 判断错误链里有没有目标错误，用法同标准库 errors.Is。
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// InvalidArgument 表示参数不合法，例如格式错误或取值不符合约定。
func InvalidArgument(msg string) Error {
	return New(CodeInvalidArgument, msg)
}

// MissingArgument 表示缺少必填参数。
func MissingArgument(msg string) Error {
	return New(CodeMissingArgument, msg)
}

// OutOfRange 表示参数格式正确，但取值超出允许范围。
func OutOfRange(msg string) Error {
	return New(CodeOutOfRange, msg)
}

// NotFound 表示按标识查询的资源不存在。
func NotFound(msg string) Error {
	return New(CodeNotFound, msg)
}

// AlreadyExists 表示资源已存在，通常由唯一性约束触发。
func AlreadyExists(msg string) Error {
	return New(CodeAlreadyExists, msg)
}

// Unauthorized 表示未通过身份认证，缺少凭证或凭证无效。
func Unauthorized(msg string) Error {
	return New(CodeUnauthorized, msg)
}

// Forbidden 表示身份已确认，但没有访问该资源的权限。
func Forbidden(msg string) Error {
	return New(CodeForbidden, msg)
}

// InternalServerError 表示服务端内部错误，调用方重试通常无法解决。
func InternalServerError(msg string) Error {
	return New(CodeInternalServerError, msg)
}

// TooManyRequests 表示请求过于频繁，调用方应降低请求速率后再试。
func TooManyRequests(msg string) Error {
	return New(CodeTooManyRequests, msg)
}

// 常用的通用错误码。
//
// code 是稳定的机器可读标识，调用方按它判断错误类型。
// 这里只放跨业务通用的语义，业务专属的错误码留在各自的领域包里。
const (
	// CodeInvalidArgument 表示参数不合法。
	CodeInvalidArgument = "invalid_argument"

	// CodeMissingArgument 表示缺少必填参数。
	CodeMissingArgument = "missing_argument"

	// CodeOutOfRange 表示参数取值超出允许范围。
	CodeOutOfRange = "out_of_range"

	// CodeNotFound 表示资源不存在。
	CodeNotFound = "not_found"

	// CodeAlreadyExists 表示资源已存在。
	CodeAlreadyExists = "already_exists"

	// CodeUnauthorized 表示未通过身份认证。
	CodeUnauthorized = "unauthorized"

	// CodeForbidden 表示已认证但缺少权限。
	CodeForbidden = "forbidden"

	// CodeTooManyRequests 表示请求过于频繁。
	CodeTooManyRequests = "too_many_requests"

	// CodeInternalServerError 表示服务端内部错误。
	CodeInternalServerError = "internal_server_error"
)
