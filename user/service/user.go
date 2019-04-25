package service

import (
	"context"
	"go-project-template/user"
	"go-project-template/user/model"
)

type UserRepository interface {
	GetUserWithId(ctx context.Context, id int) (result *model.User, err error)

	GetUserWithUsername(ctx context.Context, username string) (result *model.User, err error)

	AddUser(ctx context.Context, user *model.AddUserParam) (result *model.User, err error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (this *UserService) GetUserWithId(ctx context.Context, id int) (result *model.User, err error) {
	result, err = this.repo.GetUserWithId(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, user.UserNotExist
	}
	return result, err
}

func (this *UserService) AddUser(ctx context.Context, param *model.AddUserParam) (result *model.User, err error) {
	if param.Username == "" {
		return nil, user.UsernameExists
	}
	eUser, err := this.repo.GetUserWithUsername(ctx, param.Username)
	if err != nil {
		return nil, err
	}

	if eUser != nil {
		return nil, user.UsernameExists
	}

	return this.repo.AddUser(ctx, param)
}
