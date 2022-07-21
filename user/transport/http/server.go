package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/errors"
	"github.com/smartwalle/http4go"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/xid"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go-project-template/user/swagger"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ErrSuccess          = errors.New(100000, "成功")
	ErrInternalError    = errors.New(100001, "内部错误")
	ErrUnauthorized     = errors.New(100002, "未登录")
	ErrPermissionDenied = errors.New(100003, "没有操作权限")
)

// Response 仅用作生成 Swagger 文档使用
type Response struct {
	Code    int32       `json:"code"`              // 错误码
	Message string      `json:"message,omitempty"` // 错误消息
	Data    interface{} `json:"data,omitempty"`    // 数据
}

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
	var nLogger = log4go.New()
	nLogger.DisablePath()
	nLogger.SetPrefix("[HTTP] ")
	nLogger.AddWriter("file", log4go.NewFileWriter(log4go.LevelTrace, log4go.WithLogDir("./logs_http"), log4go.WithMaxAge(60*60*24*30)))

	var s = &Server{}
	s.engine = gin.Default()

	s.engine.Use(MidRequestTag())
	s.engine.Use(MidCORS())
	s.engine.Use(MidLogger(nLogger))

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

const (
	kRequestTag = "Request-Tag"
)

func MidRequestTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rid = xid.NewMID().Hex()
		c.Request.Header.Set(kRequestTag, rid)
		c.Writer.Header().Add(kRequestTag, rid)
	}
}

func GetRequestTag(c *gin.Context) string {
	return c.Request.Header.Get(kRequestTag)
}

func MidLogger(logger log4go.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.ParseForm()

		if logger != nil {
			var w = &strings.Builder{}
			w.WriteString(fmt.Sprintf("%s - %s \n", c.Request.Method, c.Request.URL.Path))
			w.WriteString(fmt.Sprintf("Header: \n"))
			for key, value := range c.Request.Header {
				w.WriteString(fmt.Sprintf("- %v: %v \n", key, value))
			}

			if len(c.Request.Form) > 0 {
				w.WriteString(fmt.Sprintf("Form: \n"))
				for key, value := range c.Request.Form {
					w.WriteString(fmt.Sprintf("- %v: %v \n", key, value))
				}
			}

			if c.ContentType() == "application/json" {
				var body, _ = http4go.CopyBody(c.Request)
				var bodyBytes, _ = ioutil.ReadAll(body)

				w.WriteString("Body: \n")
				w.Write(bodyBytes)
				w.WriteString("\n")
			}

			logger.Log(w.String())

		}
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
