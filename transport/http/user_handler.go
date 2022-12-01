package http

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-project-template/pkg"
	"go-project-template/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	var h = &UserHandler{}
	h.userService = userService
	return h
}

func (this *UserHandler) Handle(r gin.IRouter) {
	r.GET("/profile", pkg.JSONWrapper(this.Profile))

	r.GET("/user", pkg.JSONWrapper(this.GetUser))

	r.POST("/user", pkg.JSONWrapper(this.AddUser))
}

// Profile 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Schemes
// @Description 获取当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=UserInfo}
// @Security ApiKeyAuth
// @Router /profile [get]
func (this *UserHandler) Profile(c *gin.Context) (interface{}, error) {
	var nUser = &UserInfo{}
	nUser.Id = 1
	nUser.Username = "SmartWalle"
	nUser.FirstName = "Feng"
	nUser.LastName = "Yang"
	return nUser, nil
}

// GetUser 获取指定用户信息
// @Summary 获取指定用户信息
// @Schemes
// @Description 获取指定用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "用户 id"
// @Success 200 {object} Response{data=UserInfo}
// @Router /user [get]
func (this *UserHandler) GetUser(c *gin.Context) (interface{}, error) {
	var id = conv4go.Int64(c.Request.FormValue("id"))

	var user, err = this.userService.GetUserWithId(id)
	if err != nil {
		return nil, err
	}
	return NewUserInfo(user), nil
}

// AddUser 添加用户信息
// @Summary 添加用户信息
// @Schemes
// @Description 添加用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param object body AddUserReq true "用户信息"
// @Success 200 {object} Response{data=UserInfo}
// @Router /user [post]
func (this *UserHandler) AddUser(c *gin.Context) (interface{}, error) {
	var req *AddUserReq
	if err := pkg.BindJSON(c, &req); err != nil {
		return nil, err
	}

	var user, err = this.userService.AddUser(req.AddUserOption())
	if err != nil {
		return nil, err
	}
	return NewUserInfo(user), nil
}
