package user

import (
	"github.com/smartwalle/errors"
)

var (
	UsernameExists = errors.New(110001, "用户名已经存在")
	UserNotExist   = errors.New(110002, "用户信息不存在")
)

// User 数据实体
type User struct {
	Id        int64  `sql:"id"`         // id
	Username  string `sql:"username"`   // 用户名
	LastName  string `sql:"last_name"`  // 姓
	FirstName string `sql:"first_name"` // 名
}

// AddUserOption service 或者 repository 方法参数过多时，可以考虑使用结构体组织
type AddUserOption struct {
	Username  string // 用户名
	LastName  string // 姓
	FirstName string // 名
}
