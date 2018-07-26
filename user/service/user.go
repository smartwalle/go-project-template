package service

import (
	"go-projet-template/user"
	"github.com/smartwalle/errors"
)

var (
	UserNotExist = errors.New("200000", "用户信息不存在")
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
		return nil, UserNotExist
	}
	return result, err
}
