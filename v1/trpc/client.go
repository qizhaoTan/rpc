package trpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
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

func (c *Client) Invoke(ctx context.Context, method string, args any, reply any) error {
	if method == "" {
		return errors.New("空方法名")
	}

	if args == nil {
		return errors.New("空请求")
	}

	if reflect.ValueOf(args).IsNil() {
		return errors.New("空请求")
	}

	data, _ := json.Marshal(args)
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
