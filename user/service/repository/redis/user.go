package redis

import (
	"context"
	"fmt"
	"github.com/smartwalle/dbr"
	"go-project-template/user/model"
	"go-project-template/user/service"
)

type userRepository struct {
	service.UserRepository
	rPool *dbr.Pool
}

func NewUserRepository(rPool *dbr.Pool, repo service.UserRepository) service.UserRepository {
	var r = &userRepository{}
	r.rPool = rPool
	r.UserRepository = repo
	return r
}

func (this *userRepository) GetUserWithId(ctx context.Context, id int64) (result *model.User, err error) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	var key = fmt.Sprintf("user_%d", id)
	if err = rSess.GET(key).UnmarshalJSON(&result); err != nil || result == nil {
		result, err = this.UserRepository.GetUserWithId(ctx, id)
		if err != nil {
			return nil, err
		}
		if result != nil {
			rSess.MarshalJSONEx(key, 1800, result)
		}
	}
	return result, err
}
