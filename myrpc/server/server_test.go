package main

import (
	"context"
	"myrpc/pb"
	"testing"

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
			addr:    "localhost:50051",
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
			addr:    "localhost:50051",
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
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterHelloServer(tt.server, tt.service)
		})
	}
}
