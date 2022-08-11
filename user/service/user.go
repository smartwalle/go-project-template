package service

import (
	"fmt"
	"github.com/smartwalle/dbs"
	"go-project-template/user/model"
	"time"
)

type UserRepository interface {
	BeginTx() (dbs.TX, UserRepository)

	WithTx(tx dbs.TX) UserRepository

	GetUserWithId(id int64) (result *model.User, err error)

	GetUserWithUsername(username string) (result *model.User, err error)

	AddUser(opt AddUserOption) (result int64, err error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (this *UserService) GetUserWithId(id int64) (result *model.User, err error) {
	result = &model.User{}
	result.Id = id
	result.Username = fmt.Sprintf("rsp-%d", id)
	result.FirstName = fmt.Sprintf("first name-%d", id)
	result.LastName = fmt.Sprintf("last name-%d", id)
	return result, nil
	//result, err = this.repo.GetUserWithId(id)
	//if err != nil {
	//	return nil, err
	//}
	//if result == nil {
	//	return nil, model.UserNotExist
	//}
	//return result, err
}

func (this *UserService) AddUser(opt AddUserOption) (result *model.User, err error) {
	result = &model.User{}
	result.Id = time.Now().Unix()
	result.Username = opt.Username
	result.FirstName = opt.FirstName
	result.LastName = opt.LastName
	return result, nil

	//if opt.Username == "" {
	//	return nil, model.UsernameExists
	//}
	//
	//var tx, nUserRepo = this.repo.BeginTx()
	//defer func() {
	//	if err != nil {
	//		tx.Rollback()
	//	}
	//}()
	//
	//eUser, err := nUserRepo.GetUserWithUsername(opt.Username)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if eUser != nil {
	//	return nil, model.UsernameExists
	//}
	//
	//userId, err := nUserRepo.AddUser(opt)
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
