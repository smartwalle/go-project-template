package http

import "go-project-template/user"

// UserRsp HTTP 接口返回数据
type UserRsp struct {
	Id        int64  `json:"id"`         // id
	Username  string `json:"username"`   // 用户名
	LastName  string `json:"last_name"`  // 姓
	FirstName string `json:"first_name"` // 名
}

// NewUserRsp 转换一次，将数据实体和返回数据职责分离
func NewUserRsp(user *user.User) UserRsp {
	var rsp = UserRsp{}
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
