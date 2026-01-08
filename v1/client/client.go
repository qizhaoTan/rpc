package main

import (
	"context"
	"log"
	"os"
	"time"
	"v1/pb"
	"v1/trpc"
)

func main() {
	// 连接到 gRPC 服务器
	client, err := trpc.NewClient("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	defer client.Close()

	// 创建 Hello 客户端
	c := pb.NewHelloClient(client)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用 Hello 方法
	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.Hello(ctx, &pb.ApplyHello{Name: name})
	if err != nil {
		log.Fatalf("调用 Hello 方法失败: %v", err)
	}

	log.Printf("服务器响应: %s", r.Msg)
}
