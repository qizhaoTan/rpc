package main

import (
	"context"
	"log"
	"v2/pb"
	"v2/trpc"
)

var users = map[int64]*pb.User{
	1: {
		Uid:  1,
		Name: "Tan",
		Age:  19,
		Sex:  1,
	},
	2: {
		Uid:  2,
		Name: "Liu",
		Age:  18,
		Sex:  2,
	},
}

// server 实现 Hello 服务接口
type server struct {
	pb.IHelloService
	pb.IUserService
}

// Hello 实现 Hello 方法
func (s *server) Hello(ctx context.Context, in *pb.ApplyHello) (*pb.ReplyHello, error) {
	log.Printf("收到请求: %v", in.Name)
	return &pb.ReplyHello{Msg: "Hello, " + in.Name + "!"}, nil
}

// User 实现 User 方法
func (s *server) User(ctx context.Context, in *pb.ApplyUser) (*pb.ReplyUser, error) {
	log.Printf("收到请求: %v", in.Uid)
	return &pb.ReplyUser{User: users[in.Uid]}, nil
}

func main() {
	// 创建 gRPC 服务器
	s, err := trpc.NewServer("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听端口: %v", err)
	}

	// 注册 Hello 服务
	pb.RegisterHelloServer(s, &server{})
	// 注册 User 服务
	pb.RegisterUserServer(s, &server{})

	log.Println("gRPC 服务器启动在 :50051")
	if err := s.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
