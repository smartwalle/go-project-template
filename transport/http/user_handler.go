package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/nconv"
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

func (h *UserHandler) Handle(r gin.IRouter) {
	r.GET("/profile", pkg.JSONWrapper(h.Profile))

	r.GET("/user", pkg.JSONWrapper(h.GetUser))

	r.POST("/user", pkg.JSONWrapper(h.AddUser))
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
func (h *UserHandler) Profile(c *gin.Context) (interface{}, error) {
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
func (h *UserHandler) GetUser(c *gin.Context) (interface{}, error) {
	var id = nconv.Int64(c.Request.FormValue("id"))

	var user, err = h.userService.GetUserWithId(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return ParseUserInfo(user), nil
}

// AddUser 添加用户信息
// @Summary 添加用户信息
// @Schemes
// @Description 添加用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param object body AddUserForm true "用户信息"
// @Success 200 {object} Response{data=UserInfo}
// @Router /user [post]
func (h *UserHandler) AddUser(c *gin.Context) (interface{}, error) {
	var form *AddUserForm
	if err := pkg.BindJSON(c, &form); err != nil {
		return nil, err
	}

	var user, err = h.userService.AddUser(context.Background(), form.AddUserOptions())
	if err != nil {
		return nil, err
	}
	return ParseUserInfo(user), nil
}
