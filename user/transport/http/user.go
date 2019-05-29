package http

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-project-template/user/model"
	"go-project-template/user/service"
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
	r.GET("/user", this.GetUser)

	r.POST("/user", this.AddUser)
}

func (this *UserHandler) GetUser(c *gin.Context) {
	result, err := this.userService.GetUserWithId(c, conv4go.Int64(c.Request.FormValue("id")))
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, result)
}

func (this *UserHandler) AddUser(c *gin.Context) {
	var param = &model.AddUserParam{}
	param.Username = c.Request.FormValue("username")
	param.FirstName = c.Request.FormValue("first_name")
	param.LastName = c.Request.FormValue("last_name")

	result, err := this.userService.AddUser(c, param)
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, result)
}
