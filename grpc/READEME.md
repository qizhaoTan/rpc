# gRPC 简易示例

这是一个最简单的 gRPC 调用案例，实现了 Server 和 Client 的通信。

## 项目结构

```
grpc/
├── proto/
│   ├── hello.proto      # Protocol Buffers 定义文件
│   ├── hello.pb.go      # 生成的 Protobuf 代码
│   └── hello_grpc.pb.go # 生成的 gRPC 代码
├── server/
│   └── server.go        # gRPC 服务器实现
├── client/
│   └── client.go        # gRPC 客户端实现
├── Makefile             # 构建和运行脚本
└── grpc.md              # 本文档
```

## 功能说明

- **Hello 服务**: 提供一个简单的 RPC 方法 `Hello`
- **请求**: 接收一个名字字符串
- **响应**: 返回 "Hello, {name}!" 的问候消息

## 使用方法

### 1. 生成 Protocol Buffers 代码（如果修改了 .proto 文件）

```bash
make proto
```

### 2. 启动服务器

在 grpc 目录下运行：

```bash
make server
```

服务器将在 `localhost:50051` 端口监听。

### 3. 运行客户端

打开新的终端窗口，在 grpc 目录下运行：

```bash
make client
```

客户端会发送默认名字 "World"，你也可以指定自定义名字：

```bash
go run client/client.go "张三"
```

## 预期输出

**服务端日志：**
```
2025/01/08 xx:xx:xx gRPC 服务器启动在 :50051
2025/01/08 xx:xx:xx 收到请求: World
```

**客户端日志：**
```
2025/01/08 xx:xx:xx 服务器响应: Hello, World!
```

## 技术栈

- Go
- gRPC
- Protocol Buffers (protobuf)
