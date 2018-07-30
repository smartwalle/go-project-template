package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-projet-template/user"
)

type UserHandler struct {
	UserService user.UserService
}

func NewUserHandler(uServ user.UserService) *UserHandler {
	var h = &UserHandler{}
	h.UserService = uServ
	return h
}

func (this *UserHandler) Run(r gin.IRouter) {
	r.GET("/user", this.GetUser)
}

func (this *UserHandler) GetUser(c *gin.Context) {
	var user, err = this.UserService.User(conv4go.Int(c.Request.FormValue("id")))
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, user)
}
