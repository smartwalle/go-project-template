package grpc

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	server *grpc.Server
}

type Handler interface {
	Handle(r *grpc.Server)
}

func NewServer() *Server {
	var s = &Server{}
	s.server = grpc.NewServer()
	return s
}

func (this *Server) Run() error {
	go func() {
		listener, err := net.Listen("tcp", ":8889")
		if err != nil {
			panic(err)
		}

		if err := this.server.Serve(listener); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (this *Server) AddHandler(h Handler) {
	h.Handle(this.server)
}
