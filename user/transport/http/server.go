package http

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go-project-template/user/docs"
	"net/http"
)

type Server struct {
	engine *gin.Engine
}

type Handler interface {
	Handle(r gin.IRouter)
}

// @title           Go Project Template
// @version         1.0
// @description     This is a sample http server.

// @contact.name   Swagger Doc
// @contact.url    https://github.com/swaggo/swag/blob/master/README_zh-CN.md

// @host      localhost:8888
// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func NewServer() *Server {
	var s = &Server{}
	s.engine = gin.Default()

	s.engine.Use(MidCORS())

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return s
}

func (this *Server) Router() gin.IRouter {
	return this.engine
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

func MidCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		var header = c.Writer.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS")
		header.Set("Access-Control-Allow-Headers", "Sec-Websocket-Key, Connection, Sec-Websocket-Version, Sec-Websocket-Extensions, Upgrade, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
