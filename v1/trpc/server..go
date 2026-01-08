package trpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"v1/pb"
)

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
