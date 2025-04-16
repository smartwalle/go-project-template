package pkg

import "fmt"

type Error interface {
	Code() int

	Message() string
}

type BussinessError struct {
	code    ErrorCode
	message string
}

func NewError(code ErrorCode, message string) BussinessError {
	return BussinessError{code: code, message: message}
}

func (err BussinessError) Error() string {
	return fmt.Sprintf("%d-%s", err.code, err.message)
}

func (err BussinessError) Code() int {
	return int(err.code)
}

func (err BussinessError) Message() string {
	return err.message
}

type ErrorCode int

const (
	CodeSuccess         ErrorCode = 1000000 // 成功
	CodeInternalError   ErrorCode = 1000001 // 内部错误
	CodeOperationFailed ErrorCode = 1000002 // 操作失败
)
