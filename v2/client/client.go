package main

import (
	"context"
	"log"
	"time"
	"v2/pb"
	"v2/trpc"
)

func main() {
	// 连接到 gRPC 服务器
	client, err := trpc.NewClient("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	defer client.Close()

	user := getUser(client, 1)
	sayHello(client, user.Name)
}

func getUser(client *trpc.Client, uid int64) *pb.User {
	// 创建 User 客户端
	c := pb.NewUserClient(client)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.User(ctx, &pb.ApplyUser{Uid: uid})
	if err != nil {
		log.Fatalf("调用 User 方法失败: %v", err)
	}

	log.Printf("服务器响应: %+v", r.User)
	return r.User
}

func sayHello(client *trpc.Client, name string) {
	// 创建 Hello 客户端
	c := pb.NewHelloClient(client)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Hello(ctx, &pb.ApplyHello{Name: name})
	if err != nil {
		log.Fatalf("调用 Hello 方法失败: %v", err)
	}

	log.Printf("服务器响应: %s", r.Msg)
}
