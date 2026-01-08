package pb

import "context"

// ClientConnInterface defines the functions clients need to perform unary and
// streaming RPCs.  It is implemented by *ClientConn, and is only intended to
// be referenced by generated code.
type ClientConnInterface interface {
	// Invoke performs a unary RPC and returns after the response is received
	// into reply.
	Invoke(ctx context.Context, method string, args any, reply any) error
}

type ApplyHello struct {
	Name string
}

type ReplyHello struct {
	Msg string
}

type HelloClient struct {
	Hello func(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)
}

func NewHelloClient(c ClientConnInterface) *HelloClient {
	helloClient := &HelloClient{
		Hello: func(ctx context.Context, apply *ApplyHello) (*ReplyHello, error) {
			var reply ReplyHello
			if err := c.Invoke(ctx, "Hello", apply, &reply); err != nil {
				return nil, err
			}
			return &reply, nil
		},
	}
	return helloClient
}

type IHelloService interface {
	Hello(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)
}
