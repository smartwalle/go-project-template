package usecase

import (
	"context"
	"errors"
	"project-template/internal/domain"
	"project-template/internal/domain/vo"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	var userCase = &UserUsecase{}
	userCase.userRepo = userRepo
	return userCase
}

// DescribeProfile 获取当前登录用户信息
func (u *UserUsecase) DescribeProfile(ctx context.Context, req vo.DescribeProfileRequest) (vo.DescribeProfileResponse, error) {
	return vo.DescribeProfileResponse{}, errors.New("not found")
}
