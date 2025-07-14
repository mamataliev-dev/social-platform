package middleware_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
)

func TestValidationInterceptor_Valid(t *testing.T) {
	interceptor := middleware.ValidationInterceptor()

	req := &userauthpb.RegisterRequest{
		Email:     "valid@example.com",
		Password:  "strongpass",
		UserName:  "username",
		Nickname:  "nickname",
		Bio:       "short bio",
		AvatarUrl: "https://example.com/avatar.jpg",
	}

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{FullMethod: "/user_auth.AuthService/Register"}

	resp, err := interceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestValidationInterceptor_Invalid(t *testing.T) {
	interceptor := middleware.ValidationInterceptor()

	req := &userauthpb.RegisterRequest{
		UserName:  "!-+",
		Email:     "invalid-email",
		Password:  "123",
		Nickname:  "@@@",
		Bio:       strings.Repeat("x", 300),
		AvatarUrl: "not-a-url",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), req, &grpc.UnaryServerInfo{
		FullMethod: "/user_auth.AuthService/Register",
	}, handler)

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	assert.Contains(t, st.Message(), "invalid argument")
}
