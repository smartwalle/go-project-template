package service

import (
	"go-projet-template/user"
	"titan/app/msg"
)

type UserRepository interface {
	User(id int) (*user.User, error)
}

type userService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) user.UserService {
	return &userService{Repo: repo}
}

func (this *userService) User(id int) (result *user.User, err error) {
	result, err = this.Repo.User(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, msg.UserNotExist.Location()
	}
	return result, err
}
