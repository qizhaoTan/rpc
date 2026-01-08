package main

import (
	"errors"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewClient(network, targetAddr string) (*Client, error) {
	if network != "tcp" {
		return nil, errors.New("不支持的协议")
	}

	if targetAddr == "" {
		return nil, errors.New("空地址")
	}

	conn, err := net.Dial(network, targetAddr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn: conn,
	}
	return c, nil
}

func main() {
	// 连接到 gRPC 服务器
	client, err := NewClient("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	defer client.Close()
}
