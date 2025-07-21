// Package utils provides utility functions for mapping domain errors to gRPC errors.
// It supports Single Responsibility and Open/Closed principles.
package utils

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
)

func GrpcUserNotFoundError(err error) error {
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, errs.ErrInternal.Error())
	}
}
