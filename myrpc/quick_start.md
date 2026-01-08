1. 得先有个Client
2. 补充TestNewClient关于连接建立的测试
3. 模拟grpc生成的HelloClient
4. HelloClient必须要有Hello方法看，暂时不验证返回值，只要能调用就行
5. 让Client实现ClientConnInterface接口包含一个Invoke方法，类似grpc一样