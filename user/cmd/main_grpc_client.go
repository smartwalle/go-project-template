package main

import (
	"context"
	"fmt"
	userGrpc "go-project-template/user/api/grpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8889", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	cc := userGrpc.NewUserClient(conn)
	fmt.Println(cc.GetUserWithId(context.Background(), &userGrpc.GetUserReq{Id: 1}))
}
