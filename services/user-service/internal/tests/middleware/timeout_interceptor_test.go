package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
)

func TestTimeoutInterceptor_CustomTimeoutMethod_Success(t *testing.T) {
	info := &grpc.UnaryServerInfo{
		FullMethod: "/user_auth.AuthService/Login",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(100 * time.Millisecond) // well within 2s
		return "success", nil
	}

	resp, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestTimeoutInterceptor_DefaultTimeout_Success(t *testing.T) {
	info := &grpc.UnaryServerInfo{
		FullMethod: "/some/UnknownMethod",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(100 * time.Millisecond) // well within 10s
		return "default-success", nil
	}

	resp, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "default-success", resp)
}

func TestTimeoutInterceptor_RequestTimesOut(t *testing.T) {
	info := &grpc.UnaryServerInfo{
		FullMethod: "/user_auth.AuthService/Login",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(3 * time.Second) // exceed 2s timeout
		return "should-not-return", nil
	}

	resp, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
	assert.Contains(t, st.Message(), "request timed out after 2s on /user_auth.AuthService/Login")
}

func TestTimeoutInterceptor_CustomShortTimeoutExceeded(t *testing.T) {
	info := &grpc.UnaryServerInfo{
		FullMethod: "/user.UserService/GetUserByNickname", // 300ms timeout
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(500 * time.Millisecond) // exceed 300ms
		return "too-late", nil
	}

	resp, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
	assert.Contains(t, st.Message(), "request timed out after 300ms on /user.UserService/GetUserByNickname")
}
