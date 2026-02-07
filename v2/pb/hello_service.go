package pb

import (
	"context"
	"fmt"
	"v2/api"
)

type ApplyHello struct {
	Name string
}

type ReplyHello struct {
	Msg string
}

type HelloClient struct {
	Hello func(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)
}

func NewHelloClient(c api.ClientConnInterface) *HelloClient {
	helloFunc := func(ctx context.Context, apply *ApplyHello) (*ReplyHello, error) {
		var reply ReplyHello
		if err := c.Invoke(ctx, fmt.Sprintf("%s.%s", (*HelloClient).Name(nil), "Hello"), apply, &reply); err != nil {
			return nil, err
		}
		return &reply, nil
	}
	helloClient := &HelloClient{
		Hello: helloFunc,
	}
	return helloClient
}

func (s *HelloClient) Name() string {
	return "hello_service"
}

func RegisterHelloServer(server api.ServiceRegistrar, service IHelloService) {
	server.RegisterService("hello_service", service)
}

type IHelloService interface {
	Hello(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)
}
