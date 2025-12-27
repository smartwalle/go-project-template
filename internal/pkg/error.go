package pkg

import (
	"strings"
)

type Error interface {
	Code() string

	Message() string
}

type BusinessError struct {
	code    string
	message string
}

func NewError(code string, message string) BusinessError {
	return BusinessError{code: code, message: message}
}

func (err BusinessError) Error() string {
	return strings.Join([]string{err.code, err.message}, "-")
}

func (err BusinessError) Code() string {
	return err.code
}

func (err BusinessError) Message() string {
	return err.message
}
