package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"myrpc/pb"
	"net"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Invoke(ctx context.Context, method string, args any, reply any) error {
	if method == "" {
		return errors.New("空方法名")
	}

	if args == nil {
		return errors.New("空请求")
	}

	apply := args.(*pb.ApplyHello)
	data, _ := json.Marshal(apply)
	if _, err := c.conn.Write(data); err != nil {
		return err
	}

	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return err
	}

	buf = buf[:n]
	json.Unmarshal(buf, reply)
	return nil
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
