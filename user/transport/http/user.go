package http

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/conv4go"
	"go-project-template/user/model"
)

func (this *Server) GetUser(c *gin.Context) {
	result, err := this.userServ.GetUserWithId(c, conv4go.Int(c.Request.FormValue("id")))
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, result)
}

func (this *Server) AddUser(c *gin.Context) {
	var param = &model.AddUserParam{}
	param.Username = c.Request.FormValue("username")
	param.FirstName = c.Request.FormValue("first_name")
	param.LastName = c.Request.FormValue("last_name")

	result, err := this.userServ.AddUser(c, param)
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, result)
}
