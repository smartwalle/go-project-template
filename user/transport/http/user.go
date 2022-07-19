package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-project-template/user"
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
	r.GET("/profile", this.Profile)

	r.GET("/user", this.GetUser)

	r.POST("/user", this.AddUser)
}

// Profile 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Schemes
// @Description 获取当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} user.User
// @Security ApiKeyAuth
// @Router /profile [get]
func (this *UserHandler) Profile(c *gin.Context) {
	fmt.Println(c.GetHeader("Authorization"))
	var user = &user.User{}
	user.Id = 1
	user.Username = "SmartWalle"
	user.FirstName = "Feng"
	user.LastName = "Yang"
	c.JSON(200, user)
}

// GetUser 获取指定用户信息
// @Summary 获取指定用户信息
// @Schemes
// @Description 获取指定用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "用户 id"
// @Success 200 {object} user.User
// @Router /user [get]
func (this *UserHandler) GetUser(c *gin.Context) {
	var id = conv4go.Int64(c.Request.FormValue("id"))

	var user = &user.User{}
	user.Id = id
	user.Username = fmt.Sprintf("user-%d", id)
	user.FirstName = fmt.Sprintf("first name-%d", id)
	user.LastName = fmt.Sprintf("last name-%d", id)
	c.JSON(200, user)

	//result, err := this.userService.GetUserWithId(conv4go.Int64(c.Request.FormValue("id")))
	//if err != nil {
	//	c.JSON(200, err)
	//	return
	//}
	//c.JSON(200, result)
}

// AddUser 添加用户信息
// @Summary 添加用户信息
// @Schemes
// @Description 添加用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param object body user.AddUserParam true "用户信息"
// @Success 200 {object} user.User
// @Router /user [post]
func (this *UserHandler) AddUser(c *gin.Context) {
	var param *user.AddUserParam
	c.ShouldBindJSON(&param)

	var user = &user.User{}
	user.Id = time.Now().Unix()
	user.Username = param.Username
	user.FirstName = param.FirstName
	user.LastName = param.LastName
	c.JSON(200, user)

	//var param = &user.AddUserParam{}
	//param.Username = c.Request.FormValue("username")
	//param.FirstName = c.Request.FormValue("first_name")
	//param.LastName = c.Request.FormValue("last_name")
	//
	//result, err := this.userService.AddUser(param)
	//if err != nil {
	//	c.JSON(200, err)
	//	return
	//}
	//c.JSON(200, result)
}
