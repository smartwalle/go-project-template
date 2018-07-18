package main

import (
	"github.com/smartwalle/dbs"
	_ "github.com/go-sql-driver/mysql"
	"go-projet-template/user/service/repository/mysql"
	"go-projet-template/user/service"
	"go-projet-template/user/transport/restful"
	"github.com/gin-gonic/gin"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(192.168.1.99:3306)/v3?parseTime=true", 30, 5)

	var us = service.NewUserService(mysql.NewUserRepository(db))

	var h = &restful.Handler{}
	h.UserService = us

	var s = gin.Default()
	h.Run(s)
	s.Run(":8888")
}
