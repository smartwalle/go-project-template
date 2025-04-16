package usecase

import "project-template/internal/domain"

type SigninUsecase struct {
	userRepo domain.UserRepository
}

func NewSigninUsecase(userRepo domain.UserRepository) *SigninUsecase {
	var signinCase = &SigninUsecase{}
	signinCase.userRepo = userRepo
	return signinCase
}
