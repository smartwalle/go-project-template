package http

import (
	"github.com/gin-gonic/gin"
	"go-project-template/user/service"
)

type Server struct {
	userServ *service.UserService
}

func NewServer(userServ *service.UserService) *Server {
	var s = &Server{}
	s.userServ = userServ
	return s
}

func (this *Server) Run() error {
	var s = gin.Default()
	this.route(s)

	go func() {
		if err := s.Run(":8888"); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (this *Server) route(r gin.IRouter) {
	r.GET("/user", this.GetUser)

	r.POST("/user", this.AddUser)
}
