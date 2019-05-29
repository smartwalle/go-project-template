package micro

import (
	"github.com/micro/go-micro"
)

type Server struct {
	server micro.Service
}

type Handler interface {
	Handle(r micro.Service)
}

func NewServer() *Server {
	var s = &Server{}
	s.server = micro.NewService(
		micro.Address(":8887"),
		micro.Name("user"),
	)
	return s
}

func (this *Server) Run() error {
	go func() {
		if err := this.server.Run(); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (this *Server) AddHandler(h Handler) {
	h.Handle(this.server)
}
