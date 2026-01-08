package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		network string
		addr    string
		wantErr bool
	}{
		{
			name:    "创建TCP客户端-成功",
			network: "tcp",
			addr:    "localhost:8080",
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
			addr:    "localhost:8080",
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
			}
		})
	}
}
