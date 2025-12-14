package repository

import (
	"github.com/smartwalle/dbs"
	"project-template/internal"
	"project-template/internal/domain"
)

type userRepository struct {
	dbs.Repositoy[domain.User]
	env *internal.Enviroment
}

func NewUserRepository(env *internal.Enviroment) domain.UserRepository {
	var repo = &userRepository{}
	repo.Repositoy = dbs.NewRepository[domain.User](env.DB)
	repo.env = env
	return repo
}
