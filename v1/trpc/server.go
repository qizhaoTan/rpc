package trpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	listener net.Listener
	services map[string]any
}

func NewServer(network, targetAddr string) (*Server, error) {
	if network != "tcp" {
		return nil, errors.New("不支持的协议")
	}

	if targetAddr == "" {
		return nil, errors.New("空地址")
	}

	// 监听端口
	listener, err := net.Listen(network, targetAddr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		services: make(map[string]any),
	}
	return server, nil
}

func (s *Server) RegisterService(serverName string, impl any) {
	s.services[serverName] = impl
}

func (s *Server) Start() error {
	if len(s.services) == 0 {
		return errors.New("没有注册Services")
	}

	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}

	// 读取请求（暂时简化）
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	buf = buf[:n]

	var a Apply
	if err := json.Unmarshal(buf, &a); err != nil {
		return err
	}

	if len(a.Args) <= 0 {
		return errors.New("没有传参数")
	}

	service, ok := s.services[a.ServiceName]
	if !ok {
		return fmt.Errorf("不存在service:%s", a.ServiceName)
	}

	serviceValue := reflect.ValueOf(service)
	method := serviceValue.MethodByName(a.MethodName)
	if !method.IsValid() {
		return fmt.Errorf("service:%s 不存在method:%s", a.ServiceName, a.MethodName)
	}

	ctx := context.Background()
	// 通过反射获取方法的第二个参数类型，New出来，并通过json.Unmarshal给参数赋值
	// 约定方法签名为：func(ctx context.Context, req *ReqType) (resp any, err error)
	methodType := method.Type()
	if methodType.NumIn() != 2 {
		return fmt.Errorf("service:%s method:%s 参数数量不正确", a.ServiceName, a.MethodName)
	}

	// 第二个参数类型（通常是 *ReqType）
	apply := reflect.New(methodType.In(1).Elem())
	if err := json.Unmarshal(a.Args, apply.Interface()); err != nil {
		return err
	}

	args := []reflect.Value{reflect.ValueOf(ctx), apply}
	results := method.Call(args)
	if len(results) != 2 {
		return fmt.Errorf("service:%s method:%s Call Failed", a.ServiceName, a.MethodName)
	}

	if err, ok := results[1].Interface().(error); ok && err != nil {
		return err
	}

	reply := results[0].Interface()
	resp, _ := json.Marshal(reply)
	conn.Write(resp)
	conn.Close()
	return err
}
