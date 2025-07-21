// Package middleware_test verifies the behavior of the ValidationInterceptor.
package middleware_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
)

// mockValidator is a simple struct that implements the validator interface for testing.
type mockValidator struct {
	shouldError bool
}

func (v *mockValidator) Validate() error {
	if v.shouldError {
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	return nil
}

// TestValidationInterceptor_ValidRequest ensures that a request that passes
// validation is handled correctly.
func TestValidationInterceptor_ValidRequest(t *testing.T) {
	// Scenario: A request with valid data passes validation.
	req := &mockValidator{shouldError: false}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := middleware.ValidationInterceptor()(context.Background(), req, &grpc.UnaryServerInfo{}, handler)
	assert.NoError(t, err)
}

// TestValidationInterceptor_InvalidRequest ensures that a request that fails
// validation returns an InvalidArgument gRPC error.
func TestValidationInterceptor_InvalidRequest(t *testing.T) {
	// Scenario: A request with invalid data fails validation.
	req := &mockValidator{shouldError: true}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	_, err := middleware.ValidationInterceptor()(context.Background(), req, &grpc.UnaryServerInfo{}, handler)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestValidationInterceptor_NoValidator ensures that a request object that does
// not implement the validator interface is handled correctly.
func TestValidationInterceptor_NoValidator(t *testing.T) {
	// Scenario: A request that does not have a Validate() method is processed normally.
	req := &userauthpb.LoginRequest{} // This request does not have a Validate method.
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := middleware.ValidationInterceptor()(context.Background(), req, &grpc.UnaryServerInfo{}, handler)
	assert.NoError(t, err)
}
