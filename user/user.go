package user

import "github.com/smartwalle/errors"

var (
	UsernameExists = errors.New("110001", "用户名已经存在")
	UserNotExist   = errors.New("110002", "用户信息不存在")
)

type User struct {
	Id        int    `json:"id"          sql:"id"`
	Username  string `json:"username"    sql:"username"`
	LastName  string `json:"last_name"   sql:"last_name"`
	FirstName string `json:"first_name"  sql:"first_name"`
}

type AddUserParam struct {
	Username  string
	LastName  string
	FirstName string
}

type UserService interface {
	GetUserWithId(id int) (result *User, err error)

	AddUser(user *AddUserParam) (result *User, err error)
}
