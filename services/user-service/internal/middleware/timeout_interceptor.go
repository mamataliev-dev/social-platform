package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var methodTimeouts = map[string]time.Duration{
	"/user_auth.AuthService/Login":        2 * time.Second,
	"/user_auth.AuthService/Register":     5 * time.Second,
	"/user.UserService/GetUserByNickname": 300 * time.Millisecond,
}

const defaultTimeout = 10 * time.Second

func TimeoutInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	timeout := defaultTimeout
	if t, ok := methodTimeouts[info.FullMethod]; ok {
		timeout = t
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{})
	var resp interface{}
	var err error

	go func() {
		resp, err = handler(ctx, req)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.DeadlineExceeded, "request timed out after %s on %s", timeout, info.FullMethod)
	case <-done:
		return resp, err
	}
}
