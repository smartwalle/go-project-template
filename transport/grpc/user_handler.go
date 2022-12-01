package grpc

import (
	"context"
	grpc2 "go-project-template/api/grpc"
	"go-project-template/service"
	"google.golang.org/grpc"
)

type UserHandler struct {
	grpc2.UnimplementedUserServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	var h = &UserHandler{}
	h.userService = userService
	return h
}

func (this *UserHandler) Handle(s *grpc.Server) {
	grpc2.RegisterUserServer(s, this)
}

func (this *UserHandler) GetUserWithId(ctx context.Context, req *grpc2.GetUserReq) (*grpc2.GetUserRsp, error) {
	user, err := this.userService.GetUserWithId(req.Id)
	if err != nil {
		return nil, err
	}
	return NewGetUserRsp(user), nil
}
