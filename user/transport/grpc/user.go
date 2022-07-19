package grpc

import (
	"context"
	userGrpc "go-project-template/user/api/grpc"
	"go-project-template/user/service"
	"google.golang.org/grpc"
)

type UserHandler struct {
	userGrpc.UnimplementedUserServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	var h = &UserHandler{}
	h.userService = userService
	return h
}

func (this *UserHandler) Handle(s *grpc.Server) {
	userGrpc.RegisterUserServer(s, this)
}

func (this *UserHandler) GetUserWithId(ctx context.Context, req *userGrpc.GetUserReq) (*userGrpc.GetUserRsp, error) {
	result, err := this.userService.GetUserWithId(req.Id)
	if err != nil {
		return nil, err
	}

	var rsp = &userGrpc.GetUserRsp{}
	rsp.Id = result.Id
	rsp.Username = result.Username
	rsp.LastName = result.LastName
	rsp.FirstName = result.FirstName

	return nil, nil
}
