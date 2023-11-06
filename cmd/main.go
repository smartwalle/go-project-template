package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/log4go"
	go_project_template "go-project-template"
	"go-project-template/config"
	"go-project-template/pkg"
	"go-project-template/service"
	"go-project-template/service/repository/postgres"
	"go-project-template/service/repository/redis"
	_ "go-project-template/swagger"
	"go-project-template/transport/grpc"
	"go-project-template/transport/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// @title           Project Name
// @version         1.0
// @description     This is a sample http server.

// @contact.name   Swagger Doc
// @contact.url    https://github.com/swaggo/swag/blob/master/README_zh-CN.md

// @schemes   http
// @host      localhost:8888
// @BasePath  /api

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// main
// 在 cmd 目录中执行本命令，生成文档 swag init --parseDependency --parseInternal -o ../swagger
func main() {
	// 使用 -v 参数以查看编译信息
	var showVersion = flag.Bool("v", false, "Build version")
	flag.Parse()
	if *showVersion == true {
		fmt.Println(go_project_template.Version())
		return
	}

	// 创建需要用到的文件目录
	os.Mkdir("./files", os.ModePerm)

	conf, err := config.LoadIni("./config.ini")
	if err != nil {
		log4go.Errorln("加载配置文件发生错误: ", err)
		os.Exit(-1)
	}

	// 初始化通用日志
	pkg.InitDefaultLog(conf.Server)

	// 初始化 SQL
	var sClient = pkg.NewSQL(conf.SQL)
	defer sClient.Close()

	// 初始化redis
	var rClient = pkg.NewRedis(conf.Redis)
	defer rClient.Close()

	var waiter = &sync.WaitGroup{}

	var userRepo = redis.NewUserRepository(rClient, postgres.NewUserRepository(sClient))
	var userService = service.NewUserService(userRepo)

	// HTTP 服务
	var hServer = pkg.NewHTTPServer(conf.HTTP)
	hServer.AddHandler(http.NewUserHandler(userService))
	hServer.Run(waiter)

	var gServer = grpc.NewServer()
	gServer.AddHandler(grpc.NewUserHandler(userService))
	gServer.Run()

	var c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
MainLoop:
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break MainLoop
		}
	}
	hServer.Stop()

	log4go.Println("PID", os.Getpid(), "等待任务结束...")
	waiter.Wait()
	log4go.Println("PID", os.Getpid(), "任务完成，程序关闭。")
}
