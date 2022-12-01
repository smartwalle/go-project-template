package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nsync/singleflight"
	"go-project-template/model"
	"go-project-template/service"
)

type userRepository struct {
	service.UserRepository
	rClient dbr.UniversalClient
	flight  *singleflight.Group
}

func NewUserRepository(rPool dbr.UniversalClient, repo service.UserRepository) service.UserRepository {
	var r = &userRepository{}
	r.rClient = rPool
	r.UserRepository = repo
	r.flight = singleflight.New()
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

func (this *userRepository) GetUserWithId(id int64) (result *model.User, err error) {
	data, err := this.flight.Do(fmt.Sprintf("user_%d", id), func(key string) (interface{}, error) {
		bytes, err := this.rClient.Get(context.Background(), key).Bytes()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		var nUser *model.User
		if len(bytes) > 0 {
			json.Unmarshal(bytes, &nUser)
			if nUser != nil {
				return nUser, nil
			}
		}

		nUser, err = this.UserRepository.GetUserWithId(id)
		if err != nil {
			return nil, err
		}
		if nUser != nil {
			bytes, _ = json.Marshal(nUser)
			this.rClient.Set(context.Background(), key, bytes, 0)
		}
		return nUser, err
	})

	if err != nil {
		return nil, err
	}
	result, _ = data.(*model.User)
	return result, err
}
