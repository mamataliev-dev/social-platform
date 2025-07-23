package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client userpb.InternalUserServiceClient
}

func NewUserServiceClient(addr string) (*UserClient, error) {
	target := fmt.Sprintf("dns:///%s", addr)

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		conn:   conn,
		client: userpb.NewInternalUserServiceClient(conn),
	}, nil
}

func (u *UserClient) FetchUserByID(userID int64) (*userpb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := u.client.FetchUserProfileByID(ctx, &userpb.FetchUserProfileByIDRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *UserClient) Close() error {
	return u.conn.Close()
}
