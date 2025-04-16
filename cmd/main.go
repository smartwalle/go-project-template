package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/bootstrap"
	"github.com/smartwalle/bootstrap/http"
	"project-template/internal"
	"project-template/internal/controller"
	"project-template/internal/repository"
	"project-template/internal/usecase"
)

//	@Title			Project Name
//	@Version		1.0
//	@Description	This is a sample http server.

//	@Contact.name	Swagger Doc
//	@Contact.url	https://github.com/swaggo/swag/blob/master/README_zh-CN.md

//	@Schemes	http
//	@Host		localhost:8888
//	@BasePath	/api

//	@SecurityDefinitions.apikey	ApiKeyAuth
//	@In							header
//	@Name						Authorization

// main
func main() {
	var env = internal.NewEnviroment()
	defer env.Clear()

	if err := env.Prepare(); err != nil {
		return
	}

	// repo
	var userRepo = repository.NewUserRepository(env)

	// usecase
	var userCase = usecase.NewUserUsecase(userRepo)

	// HTTP
	var router = gin.Default()

	// controller
	controller.NewUserContoller(userCase).Route(router)

	var app = bootstrap.New(
		bootstrap.WithServers(http.NewServer(env.Config.HTTP.Addr, router)),
	)

	fmt.Println(app.Run())
}
