package main

import (
	"context"
	"fmt"
	user_api "go-project-template/user/api/grpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8889", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	cc := user_api.NewUserClient(conn)
	fmt.Println(cc.GetUserWithId(context.Background(), &user_api.GetUserReq{Id: 1}))
}
