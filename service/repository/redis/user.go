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

func NewUserRepository(rClient dbr.UniversalClient, repo service.UserRepository) service.UserRepository {
	var r = &userRepository{}
	r.rClient = rClient
	r.UserRepository = repo
	r.flight = singleflight.NewGroup[string, *model.User]()
	return r
}

func (repo *userRepository) BeginTx(ctx context.Context) (*dbs.Tx, service.UserRepository, error) {
	var tx *dbs.Tx
	var err error
	var uRepo service.UserRepository

	tx, uRepo, err = repo.UserRepository.BeginTx(ctx)
	if err != nil {
		return nil, nil, err
	}
	var nRepo = *repo
	nRepo.UserRepository = uRepo
	return tx, &nRepo, nil
}

func (repo *userRepository) WithTx(tx *dbs.Tx) service.UserRepository {
	var nRepo = *repo
	nRepo.UserRepository = repo.UserRepository.WithTx(tx)
	return &nRepo
}

func (repo *userRepository) GetUserWithId(ctx context.Context, id int64) (*model.User, error) {
	return repo.flight.Do(fmt.Sprintf("user:%d", id), func(key string) (*model.User, error) {
		bytes, err := repo.rClient.Get(ctx, key).Bytes()
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

		user, err = repo.UserRepository.GetUserWithId(ctx, id)
		if err != nil {
			return nil, err
		}
		if user != nil {
			bytes, _ = json.Marshal(user)
			repo.rClient.Set(ctx, key, bytes, 0)
		}
		return user, err
	})
}
