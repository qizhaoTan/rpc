package main

import (
	"context"
	"log"
	"v2/pb"
	"v2/trpc"
)

// server 实现 Hello 服务接口
type server struct {
	pb.IHelloService
}

// Hello 实现 Hello 方法
func (s *server) Hello(ctx context.Context, in *pb.ApplyHello) (*pb.ReplyHello, error) {
	log.Printf("收到请求: %v", in.Name)
	return &pb.ReplyHello{Msg: "Hello, " + in.Name + "!"}, nil
}

func main() {
	// 创建 gRPC 服务器
	s, err := trpc.NewServer("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听端口: %v", err)
	}

	// 注册 Hello 服务
	pb.RegisterHelloServer(s, &server{})

	log.Println("gRPC 服务器启动在 :50051")
	if err := s.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
