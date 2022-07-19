package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-project-template/user/service"
	"time"
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
// @Success 200 {object} Response{data=UserRsp}
// @Security ApiKeyAuth
// @Router /profile [get]
func (this *UserHandler) Profile(c *gin.Context) (interface{}, error) {
	fmt.Println(c.GetHeader("Authorization"))
	var rsp = &UserRsp{}
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
// @Success 200 {object} Response{data=UserRsp}
// @Router /user [get]
func (this *UserHandler) GetUser(c *gin.Context) (interface{}, error) {
	var id = conv4go.Int64(c.Request.FormValue("id"))

	var rsp = &UserRsp{}
	rsp.Id = id
	rsp.Username = fmt.Sprintf("rsp-%d", id)
	rsp.FirstName = fmt.Sprintf("first name-%d", id)
	rsp.LastName = fmt.Sprintf("last name-%d", id)
	return rsp, nil

	//nUser, err := this.userService.GetUserWithId(conv4go.Int64(c.Request.FormValue("id")))
	//if err != nil {
	//	return nil, err
	//}
	//var rsp = &UserRsp{}
	//rsp.Id = nUser.Id
	//rsp.Username = nUser.Username
	//rsp.FirstName = nUser.FirstName
	//rsp.LastName = nUser.LastName
	//return nUser, nil
}

// AddUser 添加用户信息
// @Summary 添加用户信息
// @Schemes
// @Description 添加用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param object body AddUserReq true "用户信息"
// @Success 200 {object} Response{data=UserRsp}
// @Router /user [post]
func (this *UserHandler) AddUser(c *gin.Context) (interface{}, error) {
	var param *AddUserReq
	c.ShouldBindJSON(&param)

	var rsp = &UserRsp{}
	rsp.Id = time.Now().Unix()
	rsp.Username = param.Username
	rsp.FirstName = param.FirstName
	rsp.LastName = param.LastName
	return rsp, nil

	//var param = &user.AddUserParam{}
	//param.Username = c.Request.FormValue("username")
	//param.FirstName = c.Request.FormValue("first_name")
	//param.LastName = c.Request.FormValue("last_name")
	//
	//nUser, err := this.userService.AddUser(param)
	//if err != nil {
	//	return nil, err
	//}
	//var rsp = &UserRsp{}
	//rsp.Id = nUser.Id
	//rsp.Username = nUser.Username
	//rsp.FirstName = nUser.FirstName
	//rsp.LastName = nUser.LastName
	//return nUser, nil
}
