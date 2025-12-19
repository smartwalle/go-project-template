package tool

import (
	"strings"
)

type Error interface {
	Code() string

	Message() string
}

type BussinessError struct {
	code    string
	message string
}

func NewError(code string, message string) BussinessError {
	return BussinessError{code: code, message: message}
}

func (err BussinessError) Error() string {
	return strings.Join([]string{err.code, err.message}, "-")
}

func (err BussinessError) Code() string {
	return err.code
}

func (err BussinessError) Message() string {
	return err.message
}
