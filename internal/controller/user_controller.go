package controller

import (
	"github.com/gin-gonic/gin"
	"project-template/internal/domain"
	"project-template/internal/usecase"
	"project-template/pkg"
)

type UserController struct {
	userCase *usecase.UserUsecase
}

func NewUserContoller(userCase *usecase.UserUsecase) *UserController {
	var userContoller = &UserController{}
	userContoller.userCase = userCase
	return userContoller
}

func (u *UserController) Route(r gin.IRouter) {
	r.GET("/testx", pkg.HandlerWrapper(u.Profile))
}

// Profile 获取当前登录用户信息
//
//	@Summary	获取当前登录用户信息
//	@Schemes
//	@Description	获取当前登录用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			object	body		domain.User	true	"请求参数"
//	@Success		200		{object}	pkg.Response{data=domain.User}
//	@Security		ApiKeyAuth
//	@Router			/profile [get]
func (u *UserController) Profile(c *gin.Context) (interface{}, error) {
	return domain.User{}, pkg.NewError(1000, "ssss")
}
