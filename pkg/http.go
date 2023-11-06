package pkg

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/errors"
	"github.com/smartwalle/grace"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/nhttp"
	"github.com/smartwalle/pprof4gin"
	"github.com/smartwalle/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	kHTTPHeaderRequestTag    = "Request-Tag"
	kHTTPHeaderAuthorization = "Authorization"
)

func GetRequestTag(c *gin.Context) string {
	return c.GetHeader(kHTTPHeaderRequestTag)
}

func GetAuthorization(ctx *gin.Context) string {
	return ctx.GetHeader(kHTTPHeaderAuthorization)
}

type HTTPServer struct {
	conf   HTTPConfig
	engine *gin.Engine
}

type HTTPHandler interface {
	Handle(r gin.IRouter)
}

func NewHTTPServer(conf HTTPConfig) *HTTPServer {
	var s = &HTTPServer{}
	s.conf = conf
	s.engine = gin.Default()
	s.engine.Use(gin.Recovery())
	s.engine.Use(MidRequestTag())
	s.engine.Use(MidCORS())

	s.engine.GET(filepath.Join(conf.SwaggerPath, "/swagger/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler))

	pprof4gin.Run(conf.PPROFPath, s.engine)
	return s
}

func (server *HTTPServer) Use(middleware ...gin.HandlerFunc) {
	server.engine.Use(middleware...)
}

func (server *HTTPServer) Run(waiter *sync.WaitGroup) error {
	go server.run(server.conf.Address(), waiter)
	return nil
}

func (server *HTTPServer) run(mainAddress string, waiter *sync.WaitGroup) {
	var mainServer = &http.Server{}
	mainServer.Addr = mainAddress
	mainServer.Handler = server.engine

	if err := grace.ServeWithOptions([]*http.Server{mainServer}, grace.WithWaiter(waiter)); err != nil {
		log4go.Errorf("启动 http 服务发生错误: %s \n", err)
		server.Stop()
		os.Exit(-1)
	}
}

func (server *HTTPServer) Stop() {
}

func (server *HTTPServer) AddHandler(h HTTPHandler) {
	h.Handle(server.engine.Group("/api"))
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

func MidRequestTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rid = xid.NewMID().Hex()
		c.Request.Header.Set(kHTTPHeaderRequestTag, rid)
		c.Writer.Header().Add(kHTTPHeaderRequestTag, rid)
	}
}

//func MidLog(logger log4go.Logger) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Request.ParseForm()
//
//		if logger != nil {
//			var w = &strings.Builder{}
//			w.WriteString(fmt.Sprintf("%s - %s \n", c.Request.Method, c.Request.URL.Path))
//			w.WriteString(fmt.Sprintf("Header: \n"))
//			for key, value := range c.Request.Header {
//				w.WriteString(fmt.Sprintf("- %v: %v \n", key, value))
//			}
//
//			if len(c.Request.Form) > 0 {
//				w.WriteString(fmt.Sprintf("Form: \n"))
//				for key, value := range c.Request.Form {
//					w.WriteString(fmt.Sprintf("- %v: %v \n", key, value))
//				}
//			}
//
//			if c.ContentType() == "application/json" {
//				var body, _ = nhttp.DumpRequestBody(c.Request)
//				var bodyBytes, _ = io.ReadAll(body)
//
//				w.WriteString("Body: \n")
//				w.Write(bodyBytes)
//				w.WriteString("\n")
//			}
//			logger.Log(w.String())
//			c.Set("logger", logger)
//		}
//	}
//}
//
//func getHTTPLogger(c *gin.Context) log4go.Logger {
//	var data, ok = c.Get("logger")
//	if !ok {
//		return nil
//	}
//	return data.(log4go.Logger)
//}

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

func BindForm(c *gin.Context, result interface{}) (err error) {
	if err = c.Request.ParseForm(); err != nil {
		return err
	}
	return nhttp.Bind(c.Request.Form, result)
	//if err = nhttp.Bind(c.Request.Form, result); err != nil {
	//var logger = getHTTPLogger(c)
	//if logger != nil {
	//	logger.Warnf("[%s] %s %s 绑定 HTTP 请求参数失败: %s \n", c.Request.Method, c.Request.URL.Path, err)
	//}
	//return err
	//}
	//if err = validator.Check(result); err != nil {
	//	var logger = getHTTPLogger(c)
	//	if logger != nil {
	//		logger.Warnf("[%s] %s %s HTTP 请求参数验证失败: %s \n", c.Request.Method, c.Request.URL.Path, err)
	//	}
	//	return err
	//}
	//return nil
}

func BindJSON(c *gin.Context, result interface{}) (err error) {
	return json.NewDecoder(c.Request.Body).Decode(result)
	//body, err := io.ReadAll(c.Request.Body)
	//if err = json.Unmarshal(body, &result); err != nil {
	//var logger = getHTTPLogger(c)
	//if logger != nil {
	//	logger.Warnf("[%s] %s %s 绑定 HTTP 请求参数失败: %s \n", c.Request.Method, c.Request.URL.Path, err)
	//}
	//return err
	//}
	//if err = validator.Check(result); err != nil {
	//	var logger = getHTTPLogger(c)
	//	if logger != nil {
	//		logger.Warnf("[%s] %s %s HTTP 请求参数验证失败: %s \n", c.Request.Method, c.Request.URL.Path, err)
	//	}
	//	return err
	//}
	//return nil
}
