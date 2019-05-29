package http

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

type Handler interface {
	Handle(r gin.IRouter)
}

func NewServer() *Server {
	var s = &Server{}
	s.engine = gin.Default()
	return s
}

func (this *Server) Run() error {
	go func() {
		if err := this.engine.Run(":8888"); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (this *Server) AddHandler(h Handler) {
	h.Handle(this.engine)
}
