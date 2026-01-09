package trpc

import (
	"encoding/json"
	"strings"
)

type Apply struct {
	ServiceName string
	MethodName  string
	Args        []byte
}

func NewApply(method string, args any) []byte {
	names := strings.Split(method, ".")
	if len(names) != 2 {
		panic("method must be service.method")
	}

	serviceName := names[0]
	if serviceName == "" {
		panic("serviceName is empty")
	}

	methodName := names[1]
	if methodName == "" {
		panic("methodName is empty")
	}

	argsData, err := json.Marshal(args)
	if err != nil {
		panic(err)
	}

	apply := &Apply{
		ServiceName: serviceName,
		MethodName:  methodName,
		Args:        argsData,
	}

	data, err := json.Marshal(apply)
	if err != nil {
		panic(err)
	}

	return data
}

type Reply struct {
	Data []byte
}
