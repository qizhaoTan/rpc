package main

import (
	"myrpc/pb"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServer 创建一个简单的测试 TCP 服务器
func newTestServer(t *testing.T, addr string) net.Listener {
	listener, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	return listener
}

func TestNewClient(t *testing.T) {
	// 启动测试服务器
	server := newTestServer(t, "localhost:50051")
	defer server.Close()

	tests := []struct {
		name    string
		network string
		addr    string
		wantErr bool
	}{
		{
			name:    "成功连接到服务器",
			network: "tcp",
			addr:    server.Addr().String(), // 使用测试服务器的实际地址
			wantErr: false,
		},
		{
			name:    "连接不存在的服务器-失败",
			network: "tcp",
			addr:    "localhost:50052", // 假设这个端口没有服务
			wantErr: true,
		},
		{
			name:    "无效地址-失败",
			network: "tcp",
			addr:    "invalid-address:abc",
			wantErr: true,
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
			addr:    server.Addr().String(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.network, tt.addr)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.conn) // 验证连接已建立
			}
		})
	}
}

// createTestClient 创建测试用的 Client
func createTestClient(t *testing.T) *Client {
	// 启动测试服务器
	server, _ := net.Listen("tcp", "localhost:0")
	go func() {
		for {
			conn, _ := server.Accept()
			conn.Close()
		}
	}()

	client, err := NewClient("tcp", server.Addr().String())
	require.NoError(t, err)
	return client
}

func TestNewHelloClient(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		{
			name:   "成功创建HelloClient",
			client: createTestClient(t), // 创建测试用的 Client
		},
		{
			name:   "nil client - 也成功",
			client: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helloClient := pb.NewHelloClient(tt.client)
			require.NotNil(t, helloClient)
		})
	}
}
