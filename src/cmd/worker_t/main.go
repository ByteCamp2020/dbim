package main

import (
	pb "bdim/src/api/comet/grpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:3109", opts...)
	if err != nil {
		fmt.Println("conn fail")
	}
	defer conn.Close()
	client := pb.NewCometClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pkg := pb.Package{
		Roomid: 0,
		Body:   []byte("111"),
	}
	resp, err := client.Push(ctx, &pkg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("success")
	fmt.Println(resp)
	time.Sleep(time.Second * 10)
}
