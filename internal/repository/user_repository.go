package repository

import (
	"context"
	"project-template/internal"
	"project-template/internal/domain"
)

type userRepository struct {
	env *internal.Enviroment
}

func NewUserRepository(env *internal.Enviroment) domain.UserRepository {
	var repo = &userRepository{}
	repo.env = env
	return repo
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) (err error) {
	return err
}

func (u *userRepository) FindByID(ctx context.Context, id int64) (user *domain.User, err error) {
	return user, nil
}
