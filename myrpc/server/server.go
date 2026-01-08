package main

import (
	"errors"
	"net"
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

func main() {

}
