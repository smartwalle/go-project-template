package user

import (
	"github.com/smartwalle/errors"
)

var (
	UsernameExists = errors.New(110001, "用户名已经存在")
	UserNotExist   = errors.New(110002, "用户信息不存在")
)

type User struct {
	Id        int64  `json:"id"          sql:"id"`
	Username  string `json:"username"    sql:"username"`
	LastName  string `json:"last_name"   sql:"last_name"`
	FirstName string `json:"first_name"  sql:"first_name"`
}

type AddUserParam struct {
	Username  string `form:"username"`
	LastName  string `form:"last_name"`
	FirstName string `form:"first_name"`
}
