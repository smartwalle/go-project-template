package grpc

import (
	"context"
	user_api "go-project-template/user/api/grpc"
	"go-project-template/user/service"
	"google.golang.org/grpc"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	var h = &UserHandler{}
	h.userService = userService
	return h
}

func (this *UserHandler) Handle(s *grpc.Server) {
	user_api.RegisterUserServer(s, this)
}

func (this *UserHandler) GetUserWithId(ctx context.Context, req *user_api.GetUserReq) (*user_api.GetUserRsp, error) {
	result, err := this.userService.GetUserWithId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var rsp = &user_api.GetUserRsp{}
	rsp.Id = result.Id
	rsp.Username = result.Username
	rsp.LastName = result.LastName
	rsp.FirstName = result.FirstName

	return nil, nil
}
