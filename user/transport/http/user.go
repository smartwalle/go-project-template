package http

import (
	"go-project-template/user/model"
	"go-project-template/user/service"
)

// GetUserRsp HTTP 接口返回数据
type GetUserRsp struct {
	Id        int64  `json:"id"`         // id
	Username  string `json:"username"`   // 用户名
	LastName  string `json:"last_name"`  // 姓
	FirstName string `json:"first_name"` // 名
}

// NewGetUserRsp 转换一次，将数据实体和返回数据职责分离
func NewGetUserRsp(user *model.User) *GetUserRsp {
	if user == nil {
		return nil
	}
	var rsp = &GetUserRsp{}
	rsp.Id = user.Id
	rsp.Username = user.Username
	rsp.LastName = user.LastName
	rsp.FirstName = user.FirstName
	return rsp
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

// AddUserRsp HTTP 接口返回数据
type AddUserRsp struct {
	Id        int64  `json:"id"`         // id
	Username  string `json:"username"`   // 用户名
	LastName  string `json:"last_name"`  // 姓
	FirstName string `json:"first_name"` // 名
}

// NewAddUserRsp 转换一次，将数据实体和返回数据职责分离
func NewAddUserRsp(user *model.User) *GetUserRsp {
	if user == nil {
		return nil
	}
	var rsp = &GetUserRsp{}
	rsp.Id = user.Id
	rsp.Username = user.Username
	rsp.LastName = user.LastName
	rsp.FirstName = user.FirstName
	return rsp
}
