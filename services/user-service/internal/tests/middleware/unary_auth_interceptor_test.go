// Package middleware_test verifies the behavior of the UnaryAuthInterceptor.
package middleware_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
)

const testSecret = "test-secret"

// generateJWT creates a signed JWT token for testing purposes.
func generateJWT(secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": "123",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// TestUnaryAuthInterceptor_PublicMethodSkipsAuth ensures that public methods
// (like Login and Register) bypass the authentication check.
func TestUnaryAuthInterceptor_PublicMethodSkipsAuth(t *testing.T) {
	// Scenario: A public method is called without a token.
	info := &grpc.UnaryServerInfo{FullMethod: "/user_auth.AuthService/Login"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := middleware.UnaryAuthInterceptor(context.Background(), nil, info, handler)
	assert.NoError(t, err)
}

// TestUnaryAuthInterceptor_MissingMetadata ensures that a request to a protected
// method without any metadata fails with an Unauthenticated error.
func TestUnaryAuthInterceptor_MissingMetadata(t *testing.T) {
	// Scenario: A protected method is called without any metadata.
	info := &grpc.UnaryServerInfo{FullMethod: "/user.UserService/GetUser"}
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
	// Scenario: A protected method is called with an invalid token.
	md := metadata.Pairs("authorization", "Bearer invalid-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{FullMethod: "/user.UserService/GetUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	_, err := middleware.UnaryAuthInterceptor(ctx, nil, info, handler)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

// TestUnaryAuthInterceptor_ValidToken ensures that a request with a valid JWT
// token successfully passes through the interceptor.
func TestUnaryAuthInterceptor_ValidToken(t *testing.T) {
	// Scenario: A protected method is called with a valid token.
	token, _ := generateJWT(testSecret)
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{FullMethod: "/user.UserService/GetUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := middleware.UnaryAuthInterceptor(ctx, nil, info, handler)
	assert.NoError(t, err)
}
