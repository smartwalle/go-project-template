package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-projet-template/user"
)

type Handler struct {
	UserService user.UserService
}

func (this *Handler) Run(r gin.IRouter) {
	r.GET("/user", this.GetUser)
}

func (this *Handler) GetUser(c *gin.Context) {
	var user, err = this.UserService.User(conv4go.Int(c.Request.FormValue("id")))
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, user)
}
