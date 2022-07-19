package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/log4go"
	"go-project-template/user/service"
	"go-project-template/user/service/repository/mysql"
	"go-project-template/user/service/repository/redis"
	"go-project-template/user/transport/grpc"
	"go-project-template/user/transport/http"
	"os"
	"os/signal"
	"syscall"
)

//  swag init --parseDependency --parseInternal -o ../swagger
func main() {

	// 初始化数据库日志
	var dbLogger = log4go.New()
	dbLogger.AddWriter("file", log4go.NewFileWriter(log4go.LevelTrace, log4go.WithLogDir("./logs_dbs"), log4go.WithMaxAge(60*60*24*30)))
	dbs.SetLogger(dbLogger)

	var db, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(192.168.1.99:3306)/test?parseTime=true", 30, 5)
	var rPool = dbr.NewRedis("192.168.1.99:6379", 10, 5)

	var userRepo = redis.NewUserRepository(rPool, mysql.NewUserRepository(db))
	var userService = service.NewUserService(userRepo)

	var hServer = http.NewServer()
	hServer.AddHandler(http.NewUserHandler(userService))
	hServer.Run()

	var gServer = grpc.NewServer()
	gServer.AddHandler(grpc.NewUserHandler(userService))
	gServer.Run()

	var c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
