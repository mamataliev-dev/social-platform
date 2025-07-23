package middleware_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/middleware"
)

// TestUnaryAuthInterceptor_MissingMetadata ensures that a request to a protected
// method without any metadata fails with an Unauthenticated error.
func TestUnaryAuthInterceptor_MissingMetadata(t *testing.T) {
	info := &grpc.UnaryServerInfo{FullMethod: "/chat.v1.ChatService/CreateRoom"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	_, err := middleware.UnaryAuthInterceptor(context.Background(), nil, info, handler)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

// TestUnaryAuthInterceptor_InvalidToken ensures that a request with a malformed
// or invalid JWT token fails with an Unauthenticated error.
func TestUnaryAuthInterceptor_InvalidToken(t *testing.T) {
	md := metadata.Pairs("authorization", "Bearer invalid-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{FullMethod: "/chat.v1.ChatService/CreateRoom"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	_, err := middleware.UnaryAuthInterceptor(ctx, nil, info, handler)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}
