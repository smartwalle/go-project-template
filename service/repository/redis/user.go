package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nsync/singleflight"
	"go-project-template/model"
	"go-project-template/service"
)

type userRepository struct {
	service.UserRepository
	rClient dbr.UniversalClient
	flight  singleflight.Group[string, *model.User]
}

func NewUserRepository(rPool dbr.UniversalClient, repo service.UserRepository) service.UserRepository {
	var r = &userRepository{}
	r.rClient = rPool
	r.UserRepository = repo
	r.flight = singleflight.NewGroup[string, *model.User]()
	return r
}

func (this *userRepository) BeginTx() (*dbs.Tx, service.UserRepository) {
	var repo = *this
	var tx *dbs.Tx
	tx, repo.UserRepository = this.UserRepository.BeginTx()
	return tx, &repo
}

func (this *userRepository) WithTx(tx *dbs.Tx) service.UserRepository {
	var repo = *this
	repo.UserRepository = this.UserRepository.WithTx(tx)
	return &repo
}

func (this *userRepository) GetUserWithId(id int64) (*model.User, error) {
	return this.flight.Do(fmt.Sprintf("user:%d", id), func(key string) (*model.User, error) {
		bytes, err := this.rClient.Get(context.Background(), key).Bytes()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		var user *model.User
		if len(bytes) > 0 {
			if err = json.Unmarshal(bytes, &user); err != nil {
				return nil, err
			}
			if user != nil {
				return user, nil
			}
		}

		user, err = this.UserRepository.GetUserWithId(id)
		if err != nil {
			return nil, err
		}
		if user != nil {
			bytes, _ = json.Marshal(user)
			this.rClient.Set(context.Background(), key, bytes, 0)
		}
		return user, err
	})
}
