package main

import (
	"context"
	"log"
	"net"

	pb "grpc/proto"

	"google.golang.org/grpc"
)

// server 实现 Hello 服务接口
type server struct {
	pb.UnimplementedHelloServer
}

// Hello 实现 Hello 方法
func (s *server) Hello(ctx context.Context, in *pb.ApplyHello) (*pb.ReplyHello, error) {
	log.Printf("收到请求: %v", in.GetName())
	return &pb.ReplyHello{Msg: "Hello, " + in.GetName() + "!"}, nil
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听端口: %v", err)
	}

	// 创建 gRPC 服务器
	s := grpc.NewServer()

	// 注册 Hello 服务
	pb.RegisterHelloServer(s, &server{})

	log.Println("gRPC 服务器启动在 :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
