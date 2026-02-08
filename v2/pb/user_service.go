package pb

import (
	"context"
	"fmt"
	"v2/api"
)

type User struct {
	Uid  int64
	Name string
	Age  int64
	Sex  int64
}

type ApplyUser struct {
	Uid int64
}

type ReplyUser struct {
	User *User
}

type UserClient struct {
	User func(ctx context.Context, apply *ApplyUser) (*ReplyUser, error)
}

func NewUserClient(c api.ClientConnInterface) *UserClient {
	userFunc := func(ctx context.Context, apply *ApplyUser) (*ReplyUser, error) {
		var reply ReplyUser
		if err := c.Invoke(ctx, fmt.Sprintf("%s.%s", (*UserClient).Name(nil), "User"), apply, &reply); err != nil {
			return nil, err
		}
		return &reply, nil
	}
	userClient := &UserClient{
		User: userFunc,
	}
	return userClient
}

func (s *UserClient) Name() string {
	return "user_service"
}

func RegisterUserServer(server api.ServiceRegistrar, service IUserService) {
	server.RegisterService("user_service", service)
}

type IUserService interface {
	User(ctx context.Context, apply *ApplyUser) (*ReplyUser, error)
}
