package service

import (
	"context"
	"fmt"
	"github.com/smartwalle/dbs"
	"go-project-template/model"
	"time"
)

type UserRepository interface {
	BeginTx(ctx context.Context) (*dbs.Tx, UserRepository, error)

	WithTx(tx *dbs.Tx) UserRepository

	GetUserWithId(ctx context.Context, id int64) (result *model.User, err error)

	GetUserWithUsername(ctx context.Context, username string) (result *model.User, err error)

	AddUser(ctx context.Context, opts AddUserOptions) (result int64, err error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) GetUserWithId(ctx context.Context, id int64) (result *model.User, err error) {
	result = &model.User{}
	result.Id = id
	result.Username = fmt.Sprintf("rsp-%d", id)
	result.FirstName = fmt.Sprintf("first name-%d", id)
	result.LastName = fmt.Sprintf("last name-%d", id)
	return result, nil
	//result, err = service.repo.GetUserWithId(id)
	//if err != nil {
	//	return nil, err
	//}
	//if result == nil {
	//	return nil, UserNotExist
	//}
	//return result, err
}

func (service *UserService) AddUser(ctx context.Context, opts AddUserOptions) (result *model.User, err error) {
	result = &model.User{}
	result.Id = time.Now().Unix()
	result.Username = opts.Username
	result.FirstName = opts.FirstName
	result.LastName = opts.LastName
	return result, nil

	//if opts.Username == "" {
	//	return nil, UsernameExists
	//}
	//
	//var tx, nUserRepo = service.repo.BeginTx()
	//defer func() {
	//	if err != nil {
	//		tx.Rollback()
	//	}
	//}()
	//
	//eUser, err := nUserRepo.GetUserWithUsername(opts.Username)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if eUser != nil {
	//	return nil, UsernameExists
	//}
	//
	//userId, err := nUserRepo.AddUser(opts)
	//if err != nil {
	//	return nil, err
	//}
	//
	//result, err = nUserRepo.GetUserWithId(userId)
	//if err != nil {
	//	return nil, err
	//}
	//tx.Commit()
	//
	//return result, nil
}
