package service

import (
	"github.com/gin-gonic/gin"
	"project-template/internal/domain/vo"
	"project-template/internal/usecase"
)

type UserService struct {
	userCase *usecase.UserUsecase
}

func NewUserService(userCase *usecase.UserUsecase) *UserService {
	var userContoller = &UserService{}
	userContoller.userCase = userCase
	return userContoller
}

func (u *UserService) Route(r gin.IRouter) {
	r.GET("/profile", Handle(u.DescribeProfile))
}

// DescribeProfile 获取当前登录用户信息
//
//	@Summary	获取当前登录用户信息
//	@Schemes
//	@Description	获取当前登录用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			object	body		vo.DescribeProfileRequest	true	"请求参数"
//	@Success		200		{object}	vo.DescribeProfileResponse
//	@Security		ApiKeyAuth
//	@Router			/profile [get]
func (u *UserService) DescribeProfile(c *gin.Context, req vo.DescribeProfileRequest) (vo.DescribeProfileResponse, error) {
	return u.userCase.DescribeProfile(c, req)
}
