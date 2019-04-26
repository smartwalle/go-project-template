package user

import (
	"github.com/smartwalle/errors"
)

var (
	UsernameExists = errors.New(110001, "用户名已经存在")
	UserNotExist   = errors.New(110002, "用户信息不存在")
)
