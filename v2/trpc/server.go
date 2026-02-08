package trpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			log.Printf("新连接进来 localAddr %s remoteAddr %s\n", conn.LocalAddr(), conn.RemoteAddr())
			defer func() {
				log.Printf("连接断开 localAddr %s remoteAddr %s\n", conn.LocalAddr(), conn.RemoteAddr())
				conn.Close()
			}()
			for {
				if err := s.recv(conn); err != nil {
					log.Printf("Server recv error: %v", err)
					return
				}
			}
		}()
	}
}

func (s *Server) recv(conn net.Conn) error {
	// 读取请求（暂时简化）
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	buf = buf[:n]

	var a Apply
	if err := json.Unmarshal(buf, &a); err != nil {
		return err
	}

	serviceName := a.ServiceName
	methodName := a.MethodName
	args := a.Args
	reply, err := s.call(args, serviceName, methodName)
	if err != nil {
		return err
	}

	resp, _ := json.Marshal(reply)
	conn.Write(resp)
	return nil
}

func (s *Server) call(args []byte, serviceName string, methodName string) (any, error) {
	if len(args) <= 0 {
		return nil, errors.New("没有传参数")
	}

	service, ok := s.services[serviceName]
	if !ok {
		return nil, fmt.Errorf("不存在service:%s", serviceName)
	}

	serviceValue := reflect.ValueOf(service)
	method := serviceValue.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("service:%s 不存在method:%s", serviceName, methodName)
	}

	ctx := context.Background()
	// 通过反射获取方法的第二个参数类型，New出来，并通过json.Unmarshal给参数赋值
	// 约定方法签名为：func(ctx context.Context, req *ReqType) (resp any, err error)
	methodType := method.Type()
	if methodType.NumIn() != 2 {
		return nil, fmt.Errorf("service:%s method:%s 参数数量不正确", serviceName, methodName)
	}

	// 第二个参数类型（通常是 *ReqType）
	apply := reflect.New(methodType.In(1).Elem())
	if err := json.Unmarshal(args, apply.Interface()); err != nil {
		return nil, err
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(ctx), apply})
	if len(results) != 2 {
		return nil, fmt.Errorf("service:%s method:%s Call Failed", serviceName, methodName)
	}

	if err, ok := results[1].Interface().(error); ok && err != nil {
		return nil, err
	}

	reply := results[0].Interface()
	return reply, nil
}
