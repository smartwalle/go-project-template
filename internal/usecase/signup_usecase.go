package usecase

import "project-template/internal/domain"

type SignupUsecase struct {
	userRepo domain.UserRepository
}

func NewSignupUsecase(userRepo domain.UserRepository) *SignupUsecase {
	var signupCase = &SignupUsecase{}
	signupCase.userRepo = userRepo
	return signupCase
}
