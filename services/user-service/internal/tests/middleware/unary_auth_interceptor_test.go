package middleware

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/middleware"
)

var testSecret = []byte("test_secret")

func generateJWT(secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": "1234",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func TestUnaryAuthInterceptor_MissingMetadata(t *testing.T) {
	ctx := context.Background()

	_, err := middleware.UnaryAuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, errs.ErrMissingMetadata.Error(), st.Message())
}

func TestUnaryAuthInterceptor_MissingAuthHeader(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	_, err := middleware.UnaryAuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, errs.ErrMissingAuthToken.Error(), st.Message())
}

func TestUnaryAuthInterceptor_InvalidToken(t *testing.T) {
	md := metadata.Pairs("authorization", "Bearer invalid.token.here")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := middleware.UnaryAuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, errs.ErrInvalidToken.Error(), st.Message())
}

func TestUnaryAuthInterceptor_ValidToken(t *testing.T) {
	token, err := generateJWT(testSecret)
	assert.NoError(t, err)

	os.Setenv("JWT_SECRET", string(testSecret))
	defer os.Unsetenv("JWT_SECRET")

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := middleware.UnaryAuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}
