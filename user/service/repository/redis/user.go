package redis

import (
	"fmt"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"go-project-template/user"
	"go-project-template/user/service"
)

type userRepository struct {
	service.UserRepository
	rPool dbr.Pool
}

func NewUserRepository(rPool dbr.Pool, repo service.UserRepository) service.UserRepository {
	var r = &userRepository{}
	r.rPool = rPool
	r.UserRepository = repo
	return r
}

func (this *userRepository) BeginTx() (dbs.TX, service.UserRepository) {
	var repo = *this
	var tx dbs.TX
	tx, repo.UserRepository = this.UserRepository.BeginTx()
	return tx, &repo
}

func (this *userRepository) WithTx(tx dbs.TX) service.UserRepository {
	var repo = *this
	repo.UserRepository = this.UserRepository.WithTx(tx)
	return &repo
}

func (this *userRepository) GetUserWithId(id int64) (result *user.User, err error) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	var key = fmt.Sprintf("user_%d", id)
	if err = rSess.GET(key).UnmarshalJSON(&result); err != nil || result == nil {
		result, err = this.UserRepository.GetUserWithId(id)
		if err != nil {
			return nil, err
		}
		if result != nil {
			rSess.MarshalJSONEx(key, 1800, result)
		}
	}
	return result, err
}
