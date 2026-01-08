package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	defer conn.Close()

	// 创建 Hello 客户端
	c := pb.NewHelloClient(conn)

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

	log.Printf("服务器响应: %s", r.GetMsg())
}
