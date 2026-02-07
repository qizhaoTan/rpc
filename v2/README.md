# TRPC - 简易 RPC 框架

一个基于 Go 实现的极简 RPC 框架，使用 TCP + JSON 进行通信。该项目专注于演示 RPC 的核心原理。

## 已实现功能

### 核心特性

- **TCP 通信**：基于 TCP 的网络通信
- **JSON 序列化**：使用 JSON 进行数据序列化/反序列化
- **服务注册**：支持服务动态注册到服务端
- **反射调用**：通过反射动态调用服务方法
- **并发处理**：每个连接在独立的 goroutine 中处理
- **简单协议**：Method 格式为 `service_name.method_name`

### 架构分层

```
├── api/          # 核心接口定义
│   └── api.go    # ClientConnInterface, ServiceRegistrar
├── trpc/         # RPC 框架实现
│   ├── server.go # 服务端实现
│   ├── client.go # 客户端实现
│   └── entity.go # 通信协议定义
├── pb/           # 协议定义（模拟 protobuf）
│   └── hello_service.go # Hello 服务示例
├── server/       # 服务端示例
│   └── server.go # Hello 服务实现
└── client/       # 客户端示例
    └── client.go # Hello 客户端实现
```

### 通信协议

```go
type Apply struct {
    ServiceName string  // 服务名称
    MethodName  string  // 方法名称
    Args        []byte  // JSON 编码的参数
}
```

方法签名约定：
```go
func (s *ServiceType) MethodName(ctx context.Context, req *ReqType) (*RespType, error)
```

## 快速开始

### 启动服务端

```bash
make server
# 或
go run server/server.go
```

### 启动客户端

```bash
make client
# 或
go run client/client.go [name]
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./trpc -v
```

## 添加新服务

1. 在 `pb/` 目录定义服务接口和消息类型
2. 实现客户端封装（参考 `NewHelloClient`）
3. 实现服务注册函数（参考 `RegisterHelloServer`）
4. 在 `server/` 中实现服务逻辑
5. 使用 `pb.RegisterXXXServer()` 注册服务

## 当前限制

本框架是一个最小化原型，主要用于学习 RPC 原理，不建议用于生产环境：

- ❌ 仅支持 TCP 协议
- ❌ 仅支持 JSON 序列化
- ❌ 没有连接池（每次调用创建新连接）
- ❌ 没有超时机制
- ❌ 没有消息分帧（存在粘包/半包问题）
- ❌ 没有连接复用
- ❌ 没有重试机制
- ❌ 没有 Context 传播
- ❌ 没有拦截器/中间件
- ❌ 没有服务发现和负载均衡
- ❌ 没有监控和日志
- ❌ 没有安全机制（TLS/认证）

详见 [plan.md](plan.md) 了解待实现功能清单。

## 示例

### 定义服务

```go
type HelloService interface {
    Hello(ctx context.Context, req *ApplyHello) (*ReplyHello, error)
}
```

### 实现服务

```go
type server struct{}

func (s *server) Hello(ctx context.Context, in *ApplyHello) (*ReplyHello, error) {
    return &ReplyHello{Msg: "Hello, " + in.Name + "!"}, nil
}
```

### 注册服务

```go
s, _ := trpc.NewServer("tcp", ":50051")
pb.RegisterHelloServer(s, &server{})
s.Start()
```

### 调用服务

```go
client, _ := trpc.NewClient("tcp", "localhost:50051")
c := pb.NewHelloClient(client)
resp, _ := c.Hello(context.Background(), &pb.ApplyHello{Name: "World"})
```

## 许可证

MIT
