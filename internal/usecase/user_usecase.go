package usecase

import "project-template/internal/domain"

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	var userCase = &UserUsecase{}
	userCase.userRepo = userRepo
	return userCase
}
