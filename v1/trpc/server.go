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
	services map[string]any
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
		services: make(map[string]any),
	}
	return server, nil
}

func (s *Server) RegisterService(serverName string, impl any) {
	s.services[serverName] = impl
}

func (s *Server) Start() error {
	if len(s.services) == 0 {
		return errors.New("没有注册Services")
	}

	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}

	// 读取请求（暂时简化）
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	buf = buf[:n]

	var a Apply
	json.Unmarshal(buf, &a)

	var apply pb.ApplyHello
	json.Unmarshal(a.Args, &apply)

	reply := &pb.ReplyHello{Msg: fmt.Sprintf("Hello, %s!", apply.Name)}
	resp, _ := json.Marshal(reply)
	conn.Write(resp)
	conn.Close()
	return err
}
