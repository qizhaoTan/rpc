package pb

import "context"

type ApplyHello struct {
	Name string
}

type ReplyHello struct {
	Msg string
}

type HelloClient struct {
	Hello func(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)
}

func NewHelloClient(_ any) *HelloClient {
	helloClient := &HelloClient{
		Hello: func(ctx context.Context, apply *ApplyHello) (*ReplyHello, error) { return nil, nil },
	}
	return helloClient
}
