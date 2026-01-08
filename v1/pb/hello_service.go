package pb

import (
	"context"
)

// ClientConnInterface defines the functions clients need to perform unary and
// streaming RPCs.  It is implemented by *ClientConn, and is only intended to
// be referenced by generated code.
type ClientConnInterface interface {
	// Invoke performs a unary RPC and returns after the response is received
	// into reply.
	Invoke(ctx context.Context, method string, args any, reply any) error
}

// ServiceRegistrar wraps a single method that supports service registration. It
// enables users to pass concrete types other than grpc.Server to the service
// registration methods exported by the IDL generated code.
type ServiceRegistrar interface {
	// RegisterService registers a service and its implementation to the
	// concrete type implementing this interface.  It may not be called
	// once the server has started serving.
	// desc describes the service and its methods and handlers. impl is the
	// service implementation which is passed to the method handlers.
	RegisterService(serverName string, impl any)
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

func RegisterHelloServer(server ServiceRegistrar, service IHelloService) {
	server.RegisterService("Hello", service)
}

type IHelloService interface {
	Hello(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)
}
