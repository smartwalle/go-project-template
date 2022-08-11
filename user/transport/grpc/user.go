package grpc

import (
	"go-project-template/user/api/grpc"
	"go-project-template/user/model"
)

func NewGetUserRsp(user *model.User) *grpc.GetUserRsp {
	if user == nil {
		return nil
	}

	var rsp = &grpc.GetUserRsp{}
	rsp.Id = user.Id
	rsp.Username = user.Username
	rsp.LastName = user.LastName
	rsp.FirstName = user.FirstName
	return rsp
}
