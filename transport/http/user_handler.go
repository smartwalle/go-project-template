package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
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
	r.GET("/profile", JSONWrapper(this.Profile))

	r.GET("/user", JSONWrapper(this.GetUser))

	r.POST("/user", JSONWrapper(this.AddUser))
}

// Profile 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Schemes
// @Description 获取当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=GetUserRsp}
// @Security ApiKeyAuth
// @Router /profile [get]
func (this *UserHandler) Profile(c *gin.Context) (interface{}, error) {
	fmt.Println(c.GetHeader("Authorization"))
	var rsp = &GetUserRsp{}
	rsp.Id = 1
	rsp.Username = "SmartWalle"
	rsp.FirstName = "Feng"
	rsp.LastName = "Yang"
	return rsp, nil
}

// GetUser 获取指定用户信息
// @Summary 获取指定用户信息
// @Schemes
// @Description 获取指定用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "用户 id"
// @Success 200 {object} Response{data=GetUserRsp}
// @Router /user [get]
func (this *UserHandler) GetUser(c *gin.Context) (interface{}, error) {
	var id = conv4go.Int64(c.Request.FormValue("id"))

	var user, err = this.userService.GetUserWithId(id)
	if err != nil {
		return nil, err
	}
	return NewGetUserRsp(user), nil
}

// AddUser 添加用户信息
// @Summary 添加用户信息
// @Schemes
// @Description 添加用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param object body AddUserReq true "用户信息"
// @Success 200 {object} Response{data=AddUserRsp}
// @Router /user [post]
func (this *UserHandler) AddUser(c *gin.Context) (interface{}, error) {
	var param *AddUserReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return nil, err
	}

	var user, err = this.userService.AddUser(param.AddUserOption())
	if err != nil {
		return nil, err
	}
	return NewAddUserRsp(user), nil
}
