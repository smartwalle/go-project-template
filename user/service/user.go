package service

import (
	"context"
	"github.com/smartwalle/dbs"
	"go-project-template/user"
	"go-project-template/user/model"
)

type UserRepository interface {
	BeginTx() (dbs.TX, UserRepository)

	WithTx(tx dbs.TX) UserRepository

	GetUserWithId(ctx context.Context, id int64) (result *model.User, err error)

	GetUserWithUsername(ctx context.Context, username string) (result *model.User, err error)

	AddUser(ctx context.Context, user *model.AddUserParam) (result int64, err error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (this *UserService) GetUserWithId(ctx context.Context, id int64) (result *model.User, err error) {
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

	var tx, nUserRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Close()
		}
	}()

	eUser, err := nUserRepo.GetUserWithUsername(ctx, param.Username)
	if err != nil {
		return nil, err
	}

	if eUser != nil {
		tx.Rollback()
		return nil, user.UsernameExists
	}

	userId, err := nUserRepo.AddUser(ctx, param)
	if err != nil {
		return nil, err
	}

	result, err = nUserRepo.GetUserWithId(ctx, userId)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return result, nil
}
