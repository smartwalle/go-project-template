package grpc

import (
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
	go func() {
		//if err := s.Run(":8888"); err != nil {
		//	panic(err)
		//}
	}()
	return nil
}
