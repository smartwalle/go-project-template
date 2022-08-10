package service

import (
	"github.com/smartwalle/dbs"
	"go-project-template/user"
)

type UserRepository interface {
	BeginTx() (dbs.TX, UserRepository)

	WithTx(tx dbs.TX) UserRepository

	GetUserWithId(id int64) (result *user.User, err error)

	GetUserWithUsername(username string) (result *user.User, err error)

	AddUser(opts *user.AddUserOptions) (result int64, err error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (this *UserService) GetUserWithId(id int64) (result *user.User, err error) {
	result, err = this.repo.GetUserWithId(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, user.UserNotExist
	}
	return result, err
}

func (this *UserService) AddUser(opts *user.AddUserOptions) (result *user.User, err error) {
	if opts.Username == "" {
		return nil, user.UsernameExists
	}

	var tx, nUserRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Close()
		}
	}()

	eUser, err := nUserRepo.GetUserWithUsername(opts.Username)
	if err != nil {
		return nil, err
	}

	if eUser != nil {
		tx.Rollback()
		return nil, user.UsernameExists
	}

	userId, err := nUserRepo.AddUser(opts)
	if err != nil {
		return nil, err
	}

	result, err = nUserRepo.GetUserWithId(userId)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return result, nil
}
