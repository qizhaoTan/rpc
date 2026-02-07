package trpc

import (
	"context"
	"testing"
	"time"
	"v2/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		network string
		addr    string
		wantErr bool
	}{
		{
			name:    "成功创建Server",
			network: "tcp",
			addr:    "localhost:0",
			wantErr: false,
		},
		{
			name:    "空地址-失败",
			network: "tcp",
			addr:    "",
			wantErr: true,
		},
		{
			name:    "不支持的协议-失败",
			network: "udp",
			addr:    "localhost:0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.network, tt.addr)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.NotNil(t, server.listener) // 验证 listener 已创建
			}
		})
	}
}

// serverImpl 测试用的服务实现
type serverImpl struct{}

func (s *serverImpl) Hello(ctx context.Context, apply *pb.ApplyHello) (*pb.ReplyHello, error) {
	return &pb.ReplyHello{Msg: "Hello, " + apply.Name + "!"}, nil
}

// createTestServer 创建测试服务器
func createTestServer(t *testing.T) *Server {
	server, err := NewServer("tcp", "localhost:0")
	require.NoError(t, err)
	return server
}

func TestRegisterHelloServer(t *testing.T) {
	tests := []struct {
		name      string
		server    *Server
		service   pb.IHelloService
		wantErr   bool
		wantPanic bool
		expectNil bool
	}{
		{
			name:      "成功注册Hello服务",
			server:    createTestServer(t),
			service:   &serverImpl{}, // 实现了 IHelloService 的结构
			wantErr:   false,
			expectNil: false,
		},
		{
			name:      "nil服务-失败",
			server:    createTestServer(t),
			service:   nil,
			wantErr:   true,
			expectNil: true,
		},
		{
			name:      "nil server-失败",
			server:    nil,
			service:   &serverImpl{},
			wantErr:   true,
			wantPanic: true,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if tt.wantPanic {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			}()

			pb.RegisterHelloServer(tt.server, tt.service)
		})
	}
}

func TestServer_Start(t *testing.T) {
	tests := []struct {
		name        string
		setupServer func(t *testing.T) *Server
		wantErr     bool
		testClient  bool // 是否测试客户端调用
	}{
		{
			name: "成功启动服务并处理请求",
			setupServer: func(t *testing.T) *Server {
				s := createTestServer(t)
				pb.RegisterHelloServer(s, &serverImpl{})
				return s
			},
			wantErr:    false,
			testClient: true,
		},
		{
			name: "未注册服务-启动失败",
			setupServer: func(t *testing.T) *Server {
				return createTestServer(t)
			},
			wantErr:    true,
			testClient: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer(t)

			// 在 goroutine 中启动服务
			errChan := make(chan error, 1)
			go func() {
				errChan <- server.Start()
			}()

			if tt.wantErr {
				err := <-errChan
				assert.Error(t, err)
			} else {
				// 等待服务启动
				time.Sleep(100 * time.Millisecond)

				if tt.testClient {
					// 测试客户端调用
					client, err := NewClient("tcp", server.listener.Addr().String())
					require.NoError(t, err)
					defer client.Close()

					helloClient := pb.NewHelloClient(client)
					resp, err := helloClient.Hello(context.Background(), &pb.ApplyHello{Name: "Test"})
					require.NoError(t, err)
					assert.Equal(t, "Hello, Test!", resp.Msg)
				}
			}
		})
	}
}
