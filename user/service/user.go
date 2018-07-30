package service

import (
	"github.com/smartwalle/errors"
	"go-projet-template/user"
)

var (
	UserNotExist = errors.New("200000", "用户信息不存在")
)

type UserRepository interface {
	User(id int) (*user.User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) user.UserService {
	return &userService{repo: repo}
}

func (this *userService) User(id int) (result *user.User, err error) {
	result, err = this.repo.User(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, UserNotExist
	}
	return result, err
}
