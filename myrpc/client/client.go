package main

import "errors"

type Client struct {
}

func NewClient(network, targetAddr string) (*Client, error) {
	if network != "tcp" {
		return nil, errors.New("不支持的协议")
	}

	if targetAddr == "" {
		return nil, errors.New("空地址")
	}

	return &Client{}, nil
}

func main() {

}
