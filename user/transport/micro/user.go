package micro

import (
	"context"
	"github.com/micro/go-micro"
	user_api "go-project-template/user/api/grpc"
	"go-project-template/user/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	var h = &UserHandler{}
	h.userService = userService
	return h
}

func (this *UserHandler) Handle(s micro.Service) {
	user_api.RegisterUserHandler(s.Server(), this)
}

func (this *UserHandler) GetUserWithId(ctx context.Context, req *user_api.GetUserReq, rsp *user_api.GetUserRsp) error {
	result, err := this.userService.GetUserWithId(req.Id)
	if err != nil {
		return err
	}

	rsp.Id = result.Id
	rsp.Username = result.Username
	rsp.LastName = result.LastName
	rsp.FirstName = result.FirstName

	return nil
}
