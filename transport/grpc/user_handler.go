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

func (h *UserHandler) Handle(s *grpc.Server) {
	grpc2.RegisterUserServer(s, h)
}

func (h *UserHandler) GetUserWithId(ctx context.Context, req *grpc2.GetUserReq) (*grpc2.GetUserRsp, error) {
	user, err := h.userService.GetUserWithId(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return NewGetUserRsp(user), nil
}
