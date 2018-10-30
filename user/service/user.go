package service

import (
	"go-project-template/user"
)

type UserRepository interface {
	GetUserWithId(id int) (result *user.User, err error)

	GetUserWithUsername(username string) (result *user.User, err error)

	AddUser(user *user.AddUserParam) (result *user.User, err error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) user.UserService {
	return &userService{repo: repo}
}

func (this *userService) GetUserWithId(id int) (result *user.User, err error) {
	result, err = this.repo.GetUserWithId(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, user.UserNotExist
	}
	return result, err
}

func (this *userService) AddUser(param *user.AddUserParam) (result *user.User, err error) {
	if param.Username == "" {
		return nil, user.UsernameExists
	}
	eUser, err := this.repo.GetUserWithUsername(param.Username)
	if err != nil {
		return nil, err
	}

	if eUser != nil {
		return nil, user.UsernameExists
	}

	return this.repo.AddUser(param)
}
