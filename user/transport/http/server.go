package http

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/errors"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go-project-template/user/docs"
	"net/http"
)

var (
	ErrSuccess          = errors.New(100000, "成功")
	ErrInternalError    = errors.New(100001, "内部错误")
	ErrUnauthorized     = errors.New(100002, "未登录")
	ErrPermissionDenied = errors.New(100003, "没有操作权限")
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

// @schemes   http
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

func JSON(c *gin.Context, status int, err error, data interface{}) {
	var rsp error
	if err != nil {
		switch nErr := err.(type) {
		case *errors.Error:
			rsp = nErr
		default:
			rsp = errors.New(-1, err.Error())
		}
	} else {
		rsp = ErrSuccess.WithData(data)
	}
	c.JSON(status, rsp)
}

func JSONWrapper(handler func(*gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		if handler != nil {
			result, err := handler(c)
			var status = http.StatusOK
			if err != nil {
				status = http.StatusBadRequest
			}
			JSON(c, status, err, result)
		}
	}
}
