package middleware

import (
	"context"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
)

func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	publicEndpoints := map[string]bool{
		"/user_auth.AuthService/Register":     true,
		"/user_auth.AuthService/Login":        true,
		"/user_auth.AuthService/Logout":       true,
		"/user_auth.AuthService/RefreshToken": true,
	}

	if publicEndpoints[info.FullMethod] {
		return handler(ctx, req)
	}

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, errs.ErrMissingMetadata.Error())
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, errs.ErrMissingAuthToken.Error())
	}

	tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, errs.ErrUnexpectedSigningMethod.Error())
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, errs.ErrInvalidToken.Error())
	}

	return handler(ctx, req)
}
