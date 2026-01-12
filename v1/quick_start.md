## 最简单的客户端
1. 得先有个Client
2. 补充TestNewClient关于连接建立的测试
3. 模拟grpc生成的HelloClient
4. HelloClient必须要有Hello方法看，暂时不验证返回值，只要能调用就行
5. 让Client实现ClientConnInterface接口包含一个Invoke方法，类似grpc一样
6. 补充TestHeroClient的测试用例，然后Hello方法中调用Invoke
7. 实习main.go，一个最简单的客户端实现就完成了，但是还缺少很多东西，但是的确是最简单的

## 最简单的服务端
1. 得先有个Server
2. 实现main.go，一个最简单的服务端实现就完成了

## 完善RPC框架
1. 新增Apply结构
2. 在传输的时候不再传输原结构，而是先转成Apply，再传输
3. server接收到数据，通过反射找到对应的service调用其method
4. server实现可以多次接受数据