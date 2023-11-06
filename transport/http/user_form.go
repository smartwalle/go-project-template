package http

import (
	"go-project-template/model"
	"go-project-template/service"
)

// UserInfo HTTP 接口返回数据
type UserInfo struct {
	Id        int64  `json:"id"`         // id
	Username  string `json:"username"`   // 用户名
	LastName  string `json:"last_name"`  // 姓
	FirstName string `json:"first_name"` // 名
}

// ParseUserInfo 转换一次，将数据实体和返回数据职责分离
func ParseUserInfo(user *model.User) UserInfo {
	var nUser = UserInfo{}
	nUser.Id = user.Id
	nUser.Username = user.Username
	nUser.LastName = user.LastName
	nUser.FirstName = user.FirstName
	return nUser
}

// AddUserForm HTTP 接口请求参数
type AddUserForm struct {
	Username  string `form:"username"      json:"username"`   // 用户名
	LastName  string `form:"last_name"     json:"last_name"`  // 姓
	FirstName string `form:"first_name"    json:"first_name"` // 名
}

func (form *AddUserForm) AddUserOptions() service.AddUserOptions {
	var opts = service.AddUserOptions{}
	opts.Username = form.Username
	opts.LastName = form.LastName
	opts.FirstName = form.FirstName
	return opts
}
