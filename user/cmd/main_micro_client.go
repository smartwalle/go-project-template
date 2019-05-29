package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	user_api "go-project-template/user/api/grpc"
)

func main() {
	var ms = micro.NewService()
	var cc = user_api.NewUserService("user", ms.Client())
	fmt.Println(cc.GetUserWithId(context.Background(), &user_api.GetUserReq{Id: 1}))
}
