package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myrpc/pb"
	"net"
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

func RegisterHelloServer(server *Server, service pb.IHelloService) {
}

type Server struct {
	listener net.Listener
}

func NewServer(network, targetAddr string) (*Server, error) {
	if network != "tcp" {
		return nil, errors.New("不支持的协议")
	}

	if targetAddr == "" {
		return nil, errors.New("空地址")
	}

	// 监听端口
	listener, err := net.Listen(network, targetAddr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}
	return server, nil
}

func (s *Server) Start() error {
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}

	// 读取请求（暂时简化）
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	buf = buf[:n]

	var apply pb.ApplyHello
	json.Unmarshal(buf, &apply)

	reply := &pb.ReplyHello{Msg: fmt.Sprintf("Hello, %s!", apply.Name)}
	resp, _ := json.Marshal(reply)
	conn.Write(resp)
	conn.Close()
	return err
}

func main() {
	// 创建 gRPC 服务器
	s, err := NewServer("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听端口: %v", err)
	}

	// 注册 Hello 服务
	RegisterHelloServer(s, &server{})

	log.Println("gRPC 服务器启动在 :50051")
	if err := s.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
