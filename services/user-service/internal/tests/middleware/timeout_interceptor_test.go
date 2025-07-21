// Package middleware_test verifies the behavior of the TimeoutInterceptor.
package middleware_test

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

// mockHandler is a simple gRPC handler that waits for a specified duration.
func mockHandler(ctx context.Context, duration time.Duration) (interface{}, error) {
	select {
	case <-time.After(duration):
		return "ok", nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// TestTimeoutInterceptor_CustomTimeout_Success ensures that a method with a custom
// timeout succeeds if the handler returns within the allotted time.
func TestTimeoutInterceptor_CustomTimeout_Success(t *testing.T) {
	// Scenario: A request to a method with a custom timeout (e.g., Login)
	// completes within its timeout duration (2s).
	info := &grpc.UnaryServerInfo{FullMethod: "/user_auth.AuthService/Login"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return mockHandler(ctx, 1*time.Second)
	}

	_, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)
	assert.NoError(t, err)
}

// TestTimeoutInterceptor_DefaultTimeout_Success ensures that a method without a
// custom timeout succeeds if it returns within the default timeout.
func TestTimeoutInterceptor_DefaultTimeout_Success(t *testing.T) {
	// Scenario: A request to a method without a custom timeout succeeds.
	info := &grpc.UnaryServerInfo{FullMethod: "/some.OtherService/SomeMethod"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return mockHandler(ctx, 4*time.Second)
	}

	_, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)
	assert.NoError(t, err)
}

// TestTimeoutInterceptor_RequestTimesOut ensures that a request that exceeds its
// timeout results in a DeadlineExceeded gRPC error.
func TestTimeoutInterceptor_RequestTimesOut(t *testing.T) {
	// Scenario: A request to a method with a custom timeout (2s) exceeds its duration.
	info := &grpc.UnaryServerInfo{FullMethod: "/user_auth.AuthService/Login"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return mockHandler(ctx, 3*time.Second)
	}

	_, err := middleware.TimeoutInterceptor(context.Background(), nil, info, handler)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
}
