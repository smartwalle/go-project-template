package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-project-template/user"
)

type UserHandler struct {
	userServ user.UserService
}

func NewUserHandler(uServ user.UserService) *UserHandler {
	var h = &UserHandler{}
	h.userServ = uServ
	return h
}

func (this *UserHandler) Run(r gin.IRouter) {
	r.GET("/user", this.GetUser)

	r.POST("/user", this.AddUser)
}

func (this *UserHandler) GetUser(c *gin.Context) {
	result, err := this.userServ.GetUserWithId(conv4go.Int(c.Request.FormValue("id")))
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, result)
}

func (this *UserHandler) AddUser(c *gin.Context) {
	var param = &user.AddUserParam{}
	param.Username = c.Request.FormValue("username")
	param.FirstName = c.Request.FormValue("first_name")
	param.LastName = c.Request.FormValue("last_name")

	result, err := this.userServ.AddUser(param)
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, result)
}
