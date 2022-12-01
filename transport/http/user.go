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

// NewUserInfo 转换一次，将数据实体和返回数据职责分离
func NewUserInfo(user *model.User) *UserInfo {
	if user == nil {
		return nil
	}
	var nUser = &UserInfo{}
	nUser.Id = user.Id
	nUser.Username = user.Username
	nUser.LastName = user.LastName
	nUser.FirstName = user.FirstName
	return nUser
}

// AddUserReq HTTP 接口请求参数
type AddUserReq struct {
	Username  string `form:"username"`   // 用户名
	LastName  string `form:"last_name"`  // 姓
	FirstName string `form:"first_name"` // 名
}

func (this *AddUserReq) AddUserOption() service.AddUserOptions {
	var opt = service.AddUserOptions{}
	opt.Username = this.Username
	opt.LastName = this.LastName
	opt.FirstName = this.FirstName
	return opt
}
