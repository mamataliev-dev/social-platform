// Package service_test verifies the behavior of AuthService’s logout logic.
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
)

// validLogoutRequest returns a RefreshTokenPayload with a sample token.
func validLogoutRequest() *userauthpb.RefreshTokenPayload {
	return &userauthpb.RefreshTokenPayload{
		RefreshToken: "sample-refresh-token",
	}
}

// TestLogout_Success ensures that a valid logout request deletes the refresh
// token and returns a success message.
func TestLogout_Success(t *testing.T) {
	// Scenario: A user successfully logs out, and the token is revoked.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLogoutRequest()

	mapper.On("ToGetRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	})
	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.Anything).Return(nil)
	mapper.On("ToLogoutResponse", mock.Anything).Return(&userauthpb.LogoutResponse{
		Message: "Logout successful",
	})

	resp, err := svc.Logout(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Logout successful", resp.Message)
}

// TestLogout_TokenNotFound ensures that attempting to log out with a non-existent
// token results in a NotFound gRPC error.
func TestLogout_TokenNotFound(t *testing.T) {
	// Scenario: A user attempts to log out with a token that is not in the repository.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLogoutRequest()

	mapper.On("ToGetRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	})
	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.Anything).Return(errs.ErrTokenNotFound)

	_, err := svc.Logout(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrTokenNotFound.Error(), st.Message())
}

// TestLogout_InternalError ensures that a generic storage error during token
// deletion results in an Internal gRPC error.
func TestLogout_InternalError(t *testing.T) {
	// Scenario: An unexpected internal error occurs during token revocation.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLogoutRequest()

	mapper.On("ToGetRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	})
	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.Anything).Return(errs.ErrInternal)

	_, err := svc.Logout(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}
