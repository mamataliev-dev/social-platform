package middleware_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
)

// wrapper implementing Validate for a successful request
type validRegisterRequest struct {
	*userauthpb.RegisterRequest
}

func (r *validRegisterRequest) Validate() error {
	return nil
}

type invalidRegisterRequest struct {
	*userauthpb.RegisterRequest
}

func (r *invalidRegisterRequest) Validate() error { // simulate validation errors across multiple fields
	return errors.New("email: invalid format; password: length must be >= 8")
}

func TestValidationInterceptor_Valid(t *testing.T) {
	interceptor := middleware.ValidationInterceptor()
	req := &validRegisterRequest{
		&userauthpb.RegisterRequest{
			Email:     "valid@example.com",
			Password:  "strongpass",
			Username:  "username",
			Nickname:  "nickname",
			Bio:       "short bio",
			AvatarUrl: "https://example.com/avatar.jpg",
		},
	}

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{FullMethod: "/user_auth.AuthService/Register"}
	handlerCalled := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return "ok", nil
	}

	resp, err := interceptor(ctx, req, info, handler)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
	assert.True(t, handlerCalled)
}

func TestValidationInterceptor_Invalid(t *testing.T) {
	interceptor := middleware.ValidationInterceptor()
	req := &invalidRegisterRequest{
		&userauthpb.RegisterRequest{
			Email:     "invalid-email",
			Password:  "123",
			Username:  "!-+",
			Nickname:  "@@@",
			Bio:       strings.Repeat("x", 300),
			AvatarUrl: "not-a-url",
		},
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
	// should contain our ErrInvalidArgument prefix plus underlying validation message
	assert.Contains(t, st.Message(), "invalid argument")
	assert.Contains(t, st.Message(), "email: invalid format")
	assert.Contains(t, st.Message(), "password: length must be >= 8")
}
