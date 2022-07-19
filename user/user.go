package user

import (
	"github.com/smartwalle/errors"
)

var (
	UsernameExists = errors.New(110001, "用户名已经存在")
	UserNotExist   = errors.New(110002, "用户信息不存在")
)

type User struct {
	Id        int64  `json:"id"          sql:"id"`         // id
	Username  string `json:"username"    sql:"username"`   // 用户名
	LastName  string `json:"last_name"   sql:"last_name"`  // 姓
	FirstName string `json:"first_name"  sql:"first_name"` // 名
}

type AddUserParam struct {
	Username  string `form:"username"`   // 用户名
	LastName  string `form:"last_name"`  // 姓
	FirstName string `form:"first_name"` // 名
}
