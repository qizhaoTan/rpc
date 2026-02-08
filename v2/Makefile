.PHONY: server client

server:
	go run server/server.go

client:
	go run client/client.go

# 只运行慢速的、需要环境的集成测试
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...
