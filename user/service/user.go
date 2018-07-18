package service

import (
	"titan/app/msg"
	"go-projet-template/user"
)

type IUserRepository interface {
	User(id int) (*user.User, error)
}

type UserService struct {
	Repo IUserRepository
}

func NewUserService(repo IUserRepository) user.IUserService {
	return &UserService{Repo: repo}
}

func (this *UserService) User(id int) (result *user.User, err error) {
	result, err = this.Repo.User(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, msg.UserNotExist.Location()
	}
	return result, err
}