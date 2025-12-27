package repository

import (
	"github.com/smartwalle/dbs"
	"project-template/internal/domain"
	"project-template/internal/domain/entity"
)

type userRepository struct {
	dbs.Repositoy[entity.User]
}

func NewUserRepository(db dbs.Database) domain.UserRepository {
	var repo = &userRepository{}
	repo.Repositoy = dbs.NewRepository[entity.User](db)
	return repo
}
