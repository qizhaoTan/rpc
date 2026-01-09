package trpc

import (
	"encoding/json"
	"testing"
	"v1/pb"

	"github.com/stretchr/testify/assert"
)

func TestNewApply(t *testing.T) {
	type args struct {
		method string
		args   any
	}
	tests := []struct {
		name string
		args args
		want *Apply
	}{
		{
			name: "正常情况-标准格式",
			args: args{
				method: "HelloService.Hello",
				args:   &pb.ApplyHello{Name: "Tan"},
			},
			want: &Apply{
				ServiceName: "HelloService",
				MethodName:  "Hello",
				Args:        []byte(`{"Name":"Tan"}`),
			},
		},
		{
			name: "正常情况-args为nil",
			args: args{
				method: "UserService.Login",
				args:   nil,
			},
			want: &Apply{
				ServiceName: "UserService",
				MethodName:  "Login",
				Args:        []byte("null"),
			},
		},
		{
			name: "正常情况-args为复杂嵌套结构",
			args: args{
				method: "OrderService.Create",
				args: &struct {
					ID       int
					Items    []string
					Metadata map[string]string
				}{
					ID:       123,
					Items:    []string{"item1", "item2"},
					Metadata: map[string]string{"key": "value"},
				},
			},
			want: &Apply{
				ServiceName: "OrderService",
				MethodName:  "Create",
				Args:        []byte(`{"ID":123,"Items":["item1","item2"],"Metadata":{"key":"value"}}`),
			},
		},
		{
			name: "正常情况-不同服务方法名",
			args: args{
				method: "AuthService.Login_v2",
				args:   &pb.ApplyHello{Name: "Test"},
			},
			want: &Apply{
				ServiceName: "AuthService",
				MethodName:  "Login_v2",
				Args:        []byte(`{"Name":"Test"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := json.Marshal(tt.want)
			assert.Equalf(t, data, NewApply(tt.args.method, tt.args.args), "NewApply(%v, %v)", tt.args.method, tt.args.args)
		})
	}
}

// TestNewApply_Panic 测试应该触发panic的场景
func TestNewApply_Panic(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		args     any
		panicMsg string
	}{
		{
			name:     "异常-method为空字符串",
			method:   "",
			args:     &pb.ApplyHello{Name: "Test"},
			panicMsg: "method must be service.method",
		},
		{
			name:     "异常-method只有服务名（无点号）",
			method:   "HelloService",
			args:     &pb.ApplyHello{Name: "Test"},
			panicMsg: "method must be service.method",
		},
		{
			name:     "异常-method只有方法名（无点号）",
			method:   "Hello",
			args:     &pb.ApplyHello{Name: "Test"},
			panicMsg: "method must be service.method",
		},
		{
			name:     "异常-method有多个点号",
			method:   "Service.Sub.Method",
			args:     &pb.ApplyHello{Name: "Test"},
			panicMsg: "method must be service.method",
		},
		{
			name:     "异常-method格式为Service.",
			method:   "HelloService.",
			args:     &pb.ApplyHello{Name: "Test"},
			panicMsg: "method must be service.method",
		},
		{
			name:     "异常-method格式为.Method",
			method:   ".Hello",
			args:     &pb.ApplyHello{Name: "Test"},
			panicMsg: "method must be service.method",
		},
		{
			name:     "异常-args无法序列化（channel）",
			method:   "Service.Method",
			args:     make(chan int),
			panicMsg: "",
		},
		{
			name:     "异常-args无法序列化（function）",
			method:   "Service.Method",
			args:     func() {},
			panicMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panicsf(t, func() {
				NewApply(tt.method, tt.args)
			}, "NewApply(%v, %v) should panic", tt.method, tt.args)
		})
	}
}
