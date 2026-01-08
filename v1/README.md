# tRpc v1

参考 gRPC 实现的简易 RPC 框架，v1 版本目标是实现一个最简单基础的 RPC 通信框架。

## 项目概述

tRpc 是一个用 Go 语言实现的轻量级 RPC 框架，采用 TCP 协议进行通信，使用 JSON 进行数据编码。本项目旨在通过最简实现展示 RPC
框架的核心概念和工作原理。

### 核心特性

- 基于 TCP 协议的 RPC 通信
- JSON 序列化/反序列化
- 服务注册和调用机制
- 类型安全的客户端 API
- gRPC 风格的接口设计
- 完整的单元测试和集成测试

### 技术栈

- Go 1.24.0
- TCP 协议
- JSON 编码
- testify 测试框架

## 项目结构

```
v1/
├── trpc/                 # 核心 RPC 框架实现
│   ├── server.go         # 服务端实现
│   ├── client.go         # 客户端实现
│   ├── server_test.go    # 服务端测试
│   └── client_test.go    # 客户端测试
├── pb/                   # 协议定义和接口
│   └── hello_service.go  # Hello 服务定义
├── server/               # 服务端示例程序
│   └── server.go         # 服务端主程序
├── client/               # 客户端示例程序
│   └── client.go         # 客户端主程序
├── go.mod                # 模块依赖定义
├── go.sum                # 依赖校验文件
├── README.md             # 项目说明
└── quick_start.md        # 快速开始指南
```

## 快速开始

### 环境要求

- Go 1.24.0 或更高版本

### 安装

```bash
cd v1
go mod download
```

### 运行示例

1. 启动服务端（在终端 1）：

```bash
go run server/server.go
```

服务端将监听在 `localhost:50051`

2. 启动客户端（在终端 2）：

```bash
go run client/client.go
```

客户端将调用 Hello 方法并显示响应

3. 自定义参数调用：

```bash
go run client/client.go "张三"
```

### 运行测试

```bash
# 运行所有测试
go test ./trpc/... -v

# 运行服务端测试
go test ./trpc -run TestServer -v

# 运行客户端测试
go test ./trpc -run TestClient -v

# 查看测试覆盖率
go test ./trpc/... -cover
```

## 核心概念

### 服务端 (Server)

服务端负责监听网络连接、接收请求、调用对应的服务方法并返回响应。

**主要组件：**

- `listener`: TCP 监听器，接受客户端连接
- `services`: 服务注册表，存储服务名称到实现的映射

**工作流程：**

1. 创建 TCP 监听器
2. 注册服务实现
3. 启动服务循环：
    - 接受客户端连接
    - 读取请求数据
    - 解析 JSON 请求
    - 查找并调用对应服务方法
    - 序列化响应并返回
    - 关闭连接

### 客户端 (Client)

客户端负责连接服务器、发送请求并接收响应。

**主要组件：**

- `conn`: TCP 连接对象

**工作流程：**

1. 建立到服务器的 TCP 连接
2. 序列化请求参数为 JSON
3. 发送请求到服务器
4. 接收响应数据
5. 反序列化响应

### 服务注册机制

通过接口抽象实现服务的注册和调用：

- `ServiceRegistrar`: 服务注册接口，定义 `RegisterService` 方法
- `ClientConnInterface`: 客户端连接接口，定义 `Invoke` 方法
- 服务实现者只需实现业务接口（如 `IHelloService`）即可注册

### 通信协议

当前实现采用简单的协议格式：

- 传输层：TCP
- 编码层：JSON
- 请求格式：`{"service": "服务名", "method": "方法名", "args": 参数对象}`
- 响应格式：JSON 编码的响应对象
- 缓冲区：固定 1024 字节

## API 文档

### 服务端 API

| 方法              | 签名                                                  | 说明                    |
|-----------------|-----------------------------------------------------|-----------------------|
| NewServer       | `func(network, targetAddr string) (*Server, error)` | 创建 RPC 服务器，仅支持 tcp 协议 |
| RegisterService | `func(s *Server, serverName string, impl any)`      | 注册服务实现到服务注册表          |
| Start           | `func(s *Server) error`                             | 启动服务器并处理请求（单次）        |

### 客户端 API

| 方法        | 签名                                                                                      | 说明          |
|-----------|-----------------------------------------------------------------------------------------|-------------|
| NewClient | `func(network, targetAddr string) (*Client, error)`                                     | 创建客户端并连接服务器 |
| Invoke    | `func(c *Client) Invoke(ctx context.Context, method string, args any, reply any) error` | 通用 RPC 调用方法 |
| Close     | `func(c *Client) error`                                                                 | 关闭客户端连接     |

### Hello 服务 API

| 类型/方法               | 签名                                                                                        | 说明                |
|---------------------|-------------------------------------------------------------------------------------------|-------------------|
| ApplyHello          | `struct { Name string }`                                                                  | Hello 服务请求参数      |
| ReplyHello          | `struct { Msg string }`                                                                   | Hello 服务响应结果      |
| NewHelloClient      | `func(c ClientConnInterface) *HelloClient`                                                | 创建类型安全的 Hello 客户端 |
| Hello               | `func(c *HelloClient) Hello(ctx context.Context, apply *ApplyHello) (*ReplyHello, error)` | 调用 Hello 方法       |
| RegisterHelloServer | `func(server ServiceRegistrar, service IHelloService)`                                    | 注册 Hello 服务到服务器   |

## 使用示例

### 服务端示例

见server/server.go

### 客户端示例

见client/client.go

## 测试

项目采用表格驱动测试（Table-Driven Tests），覆盖了主要功能和错误场景。

### 测试覆盖范围

**服务端测试 (server_test.go)：**

- 服务器创建（成功、失败场景）
- 服务注册（正常、nil 处理）
- 服务启动和请求处理（包含完整的客户端调用）

**客户端测试 (client_test.go)：**

- 客户端创建（成功、连接失败、无效地址）
- RPC 调用（正常、空方法名、nil 参数）
- 模拟服务器隔离测试

### 运行测试

```bash
# 运行所有测试并显示详细输出
go test ./trpc/... -v

# 查看测试覆盖率
go test ./trpc/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 设计思路

### 简单优先原则

v1 版本遵循最小可行实现原则：

- 单连接处理模式
- 固定大小缓冲区（1024 字节）
- 简单的错误处理
- 基础的 JSON 编码

### 参考 gRPC 设计

接口设计参考了 gRPC 的风格：

- `Invoke` 方法模拟 gRPC 的调用接口
- 服务注册机制（`RegisterService`）
- 类型安全的客户端封装
- 上下文（Context）支持

### 渐进式开发

项目采用渐进式开发方式，从最简单的实现开始：

1. 实现基础的 TCP 通信
2. 添加 JSON 编码
3. 实现服务注册机制
4. 封装类型安全的客户端
5. 完善测试覆盖

详细开发过程参见 [quick_start.md](quick_start.md)

## 参考资源

- [quick_start.md](quick_start.md) - 详细的开发步骤记录
- [gRPC 官方文档](https://grpc.io/docs/)
- [Go Net 包文档](https://pkg.go.dev/net)

## 许可证

本项目为学习练习项目，仅供参考。
