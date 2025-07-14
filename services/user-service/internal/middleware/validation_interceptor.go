package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
)

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if msg, ok := req.(interface{ Validate() error }); ok {
			if err := msg.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "%s: %s", errs.ErrInvalidArgument.Error(), err.Error())
			}
		}
		return handler(ctx, req)
	}
}
