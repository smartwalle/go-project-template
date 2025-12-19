package repository

import (
	"github.com/smartwalle/dbs"
	"project-template/internal"
	"project-template/internal/domain"
	"project-template/internal/domain/entity"
)

type userRepository struct {
	dbs.Repositoy[entity.User]
	env *internal.Enviroment
}

func NewUserRepository(env *internal.Enviroment) domain.UserRepository {
	var repo = &userRepository{}
	repo.Repositoy = dbs.NewRepository[entity.User](env.DB)
	repo.env = env
	return repo
}
