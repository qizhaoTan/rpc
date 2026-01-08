package main

import (
	"errors"
	"net"
)

type Client struct {
	conn net.Conn
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

}
