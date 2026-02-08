//go:build integration

package trpc

import (
	"context"
	"sync"
	"testing"
	"time"
	"v2/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServerImpl 测试用服务实现
type testServerImpl struct{}

func (s *testServerImpl) Hello(ctx context.Context, apply *pb.ApplyHello) (*pb.ReplyHello, error) {
	return &pb.ReplyHello{Msg: "Hello, " + apply.Name + "!"}, nil
}

// TestIntegrationE2E 端到端集成测试
// 启动一个真实server goroutine和一个client，验证完整通信流程
func TestIntegrationE2E(t *testing.T) {
	// 创建并启动服务器
	server, err := NewServer("tcp", "localhost:0")
	require.NoError(t, err)
	require.NotNil(t, server)

	// 注册服务
	pb.RegisterHelloServer(server, &testServerImpl{})

	// 创建用于同步的channel
	serverReady := make(chan struct{}, 1)
	serverErr := make(chan error, 1)

	// 启动服务器goroutine
	go func() {
		serverReady <- struct{}{} // 通知服务器已准备就绪
		if err := server.Start(); err != nil {
			serverErr <- err
		}
	}()

	// 等待服务器启动
	<-serverReady
	time.Sleep(50 * time.Millisecond) // 给服务器一点时间完全启动

	// 获取服务器监听的地址
	serverAddr := server.listener.Addr().String()
	t.Logf("服务器监听地址: %s", serverAddr)

	// 启动客户端goroutine
	clientErr := make(chan error, 1)
	var resp *pb.ReplyHello

	go func() {
		defer close(clientErr)

		// 创建客户端
		client, err := NewClient("tcp", serverAddr)
		if err != nil {
			clientErr <- err
			return
		}
		defer client.Close()

		// 创建 Hello 客户端
		helloClient := pb.NewHelloClient(client)

		// 调用 Hello 方法
		resp, err = helloClient.Hello(context.Background(), &pb.ApplyHello{Name: "IntegrationTest"})
		if err != nil {
			clientErr <- err
			return
		}
		clientErr <- nil
	}()

	// 等待客户端完成
	select {
	case err := <-clientErr:
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Hello, IntegrationTest!", resp.Msg)
	case <-time.After(5 * time.Second):
		t.Fatal("客户端调用超时")
	}

	t.Log("端到端测试成功")
}

// TestIntegrationConcurrent 并发集成测试
// 验证多个客户端同时调用服务时的正确性
func TestIntegrationConcurrent(t *testing.T) {
	// 创建并启动服务器
	server, err := NewServer("tcp", "localhost:0")
	require.NoError(t, err)

	// 注册服务
	pb.RegisterHelloServer(server, &testServerImpl{})

	// 启动服务器goroutine
	serverReady := make(chan struct{}, 1)
	go func() {
		serverReady <- struct{}{}
		_ = server.Start()
	}()

	// 等待服务器启动
	<-serverReady
	time.Sleep(50 * time.Millisecond)

	// 获取服务器监听的地址
	serverAddr := server.listener.Addr().String()
	t.Logf("服务器监听地址: %s", serverAddr)

	// 并发参数
	concurrency := 10
	requests := 100

	// 启动多个客户端并发调用
	var wg sync.WaitGroup
	wg.Add(concurrency)

	results := make(chan *pb.ReplyHello, requests)
	errors := make(chan error, requests)

	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < requests/concurrency; j++ {
				// 创建客户端
				client, err := NewClient("tcp", serverAddr)
				if err != nil {
					errors <- err
					continue
				}
				defer client.Close()

				// 创建 Hello 客户端
				helloClient := pb.NewHelloClient(client)

				// 调用 Hello 方法
				resp, err := helloClient.Hello(context.Background(), &pb.ApplyHello{
					Name: "ConcurrentTest",
				})
				if err != nil {
					errors <- err
					continue
				}
				results <- resp
			}
		}(i)
	}

	// 等待所有goroutine完成
	wg.Wait()
	close(results)
	close(errors)

	// 统计结果
	successCount := 0
	for resp := range results {
		successCount++
		assert.Equal(t, "Hello, ConcurrentTest!", resp.Msg)
	}

	errorCount := len(errors)
	for err := range errors {
		t.Errorf("调用失败: %v", err)
	}

	t.Logf("成功调用: %d, 失败调用: %d", successCount, errorCount)
	assert.Equal(t, requests, successCount, "所有请求都应该成功")
	assert.Equal(t, 0, errorCount, "不应该有错误")
}

// TestIntegrationSameClientMultipleCalls 同一客户端多次调用集成测试
// 验证同一个客户端对象可以进行多次RPC调用
func TestIntegrationSameClientMultipleCalls(t *testing.T) {
	// 创建并启动服务器
	server, err := NewServer("tcp", "localhost:0")
	require.NoError(t, err)

	// 注册服务
	pb.RegisterHelloServer(server, &testServerImpl{})

	// 启动服务器goroutine
	serverReady := make(chan struct{}, 1)
	go func() {
		serverReady <- struct{}{}
		_ = server.Start()
	}()

	// 等待服务器启动
	<-serverReady
	time.Sleep(50 * time.Millisecond)

	// 获取服务器监听的地址
	serverAddr := server.listener.Addr().String()
	t.Logf("服务器监听地址: %s", serverAddr)

	// 创建客户端（只创建一次）
	client, err := NewClient("tcp", serverAddr)
	require.NoError(t, err)
	defer client.Close()

	// 创建 Hello 客户端
	helloClient := pb.NewHelloClient(client)

	// 第一次调用
	t.Log("第一次调用...")
	resp1, err := helloClient.Hello(context.Background(), &pb.ApplyHello{Name: "FirstCall"})
	require.NoError(t, err)
	assert.Equal(t, "Hello, FirstCall!", resp1.Msg)
	t.Logf("第一次调用成功: %s", resp1.Msg)

	// 第二次调用（使用同一个client对象）
	t.Log("第二次调用...")
	resp2, err := helloClient.Hello(context.Background(), &pb.ApplyHello{Name: "SecondCall"})
	require.NoError(t, err)
	assert.Equal(t, "Hello, SecondCall!", resp2.Msg)
	t.Logf("第二次调用成功: %s", resp2.Msg)

	// 验证两次调用的响应是不同的
	assert.NotEqual(t, resp1.Msg, resp2.Msg, "两次调用的响应应该不同")

	t.Log("同一客户端多次调用测试成功")
}
