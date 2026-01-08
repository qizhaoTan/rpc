package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
