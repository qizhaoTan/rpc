# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个基于 Go 的自定义 RPC 框架实现（v2 版本），使用 TCP + JSON 作为通信协议。框架通过反射实现服务发现和方法调用，提供了类似 gRPC 的开发体验。

## 核心架构

### 分层设计

1. **api/** - 核心接口定义
   - `ClientConnInterface`: 客户端接口，定义 `Invoke` 方法用于 RPC 调用
   - `ServiceRegistrar`: 服务端接口，定义 `RegisterService` 方法用于服务注册

2. **trpc/** - RPC 框架核心实现
   - `Server`: 服务端，负责监听连接、接收请求、通过反射调用服务方法
   - `Client`: 客户端，负责建立连接、发送请求、接收响应
   - `entity.go`: 定义通信协议结构（Apply/Reply）

3. **pb/** - 协议定义（模拟 protobuf 生成的代码）
   - 定义请求/响应结构体
   - 实现客户端封装（如 `HelloClient`）
   - 实现服务注册函数（如 `RegisterHelloServer`）

4. **server/** 和 **client/** - 示例应用程序

### 通信协议

所有通信使用 JSON 格式，消息结构：

```go
type Apply struct {
    ServiceName string  // 服务名称
    MethodName  string  // 方法名称
    Args        []byte  // JSON 编码的参数
}
```

方法名格式为 `service_name.method_name`，如 `"hello_service.Hello"`。

### 服务端工作流程

1. `Server.RegisterService()` 注册服务实现
2. `Server.Start()` 启动监听，为每个连接启动 goroutine
3. `recv()` 读取请求，解析 JSON
4. `call()` 通过反射找到服务方法，动态创建参数并调用
5. 将响应 JSON 写回连接

### 方法签名约定

服务方法必须遵循以下签名：

```go
func (s *ServiceType) MethodName(ctx context.Context, req *ReqType) (*RespType, error)
```

- 第一个参数：`context.Context`
- 第二个参数：请求结构体指针
- 返回值：响应结构体指针和 error

## 常用命令

### 运行示例

```bash
# 启动服务端
make server

# 在另一个终端启动客户端
make client
```

或直接使用：

```bash
go run server/server.go
go run client/client.go [name]  # 可选参数
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./trpc -v

# 运行单个测试
go test ./trpc -v -run TestServer_Start
```

## 添加新服务

1. 在 `pb/` 目录下定义新的服务接口和消息类型
2. 实现服务客户端（如 `NewHelloClient`）
3. 实现服务注册函数（如 `RegisterHelloServer`）
4. 在 `server/` 中实现服务结构体
5. 使用 `pb.RegisterXXXServer()` 注册服务

## 关键特性

- **仅支持 TCP 协议**：客户端和服务端都硬编码为只支持 TCP
- **JSON 序列化**：所有数据使用 JSON 编码/解码
- **反射调用**：服务端使用反射动态调用方法，无需代码生成
- **并发处理**：每个连接在独立的 goroutine 中处理
- **简单协议**：消息格式简单，易于调试和扩展

## 当前局限性

⚠️ **这是一个最小化原型，仅用于学习 RPC 原理，不建议用于生产环境**

### 已知问题

1. **消息分帧缺失**（trpc/server.go:63-85）
   - 没有长度前缀，无法正确处理 TCP 粘包/半包
   - 固定 1024 字节缓冲区，无法处理大消息

2. **超时机制缺失**（trpc/client.go:53-54, trpc/server.go:50）
   - 所有网络操作没有超时，可能导致永久阻塞
   - Read() 可能永久等待

3. **连接管理问题**（trpc/client.go:24）
   - 每次调用创建新连接，性能开销大
   - 没有连接池和连接复用

4. **错误处理不完善**（trpc/server.go:82, trpc/client.go:60）
   - json.Marshal 和 json.Unmarshal 的错误被忽略
   - 错误信息不够详细

5. **Context 传播缺失**（trpc/server.go:104）
   - 硬编码 context.Background()，无法取消请求
   - 客户端的 context 丢失

6. **没有拦截器/中间件**
   - 无法添加日志、监控、认证等功能
   - 扩展性差

7. **没有可靠性机制**
   - 无重试、无熔断、无限流
   - 无优雅关闭

8. **没有可观测性**
   - 无监控指标
   - 无分布式追踪
   - 无结构化日志

### 待实现功能清单

详见 [plan.md](plan.md)，主要包括：

- **P0（必须）**：消息分帧、超时机制、连接复用、错误处理、Context 传播
- **P1（重要）**：拦截器、重试、优雅关闭、健康检查、流式通信
- **P2（增强）**：负载均衡、服务发现、监控、TLS、认证、代码生成

### 参考资源

- 生产级 RPC 框架：gRPC、Dubbo、brpc
- RPC 设计原理：[gRPC 官方文档](https://grpc.io/docs/)
- 代码结构参考：Go 标准库 net/rpc

## AI行为
- 每次回复最后面必须携带[rpc-v2]