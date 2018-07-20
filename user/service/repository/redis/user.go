package redis

import (
	"fmt"
	"github.com/smartwalle/dbr"
	"go-projet-template/user"
	"go-projet-template/user/service"
)

type userRepository struct {
	rPool *dbr.Pool
	repo  service.UserRepository
}

func NewUserRepository(rPool *dbr.Pool, repo service.UserRepository) service.UserRepository {
	return &userRepository{rPool: rPool, repo: repo}
}

func (this *userRepository) User(id int) (result *user.User, err error) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	var key = fmt.Sprintf("user_%d", id)

	if err = rSess.GET(key).UnmarshalJSON(&result); err != nil || result == nil {
		result, err = this.repo.User(id)
		if err != nil {
			return nil, err
		}
		if result != nil {
			rSess.MarshalJSONEx(key, 1800, result)
		}
	}
	return result, err
}
