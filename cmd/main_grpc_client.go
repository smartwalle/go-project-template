package main

import (
	"context"
	"fmt"
	grpc2 "go-project-template/api/grpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8889", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	cc := grpc2.NewUserClient(conn)
	fmt.Println(cc.GetUserWithId(context.Background(), &grpc2.GetUserReq{Id: 1}))
}
