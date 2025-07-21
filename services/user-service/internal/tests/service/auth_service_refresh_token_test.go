// Package service_test verifies the behavior of AuthService’s refresh token logic.
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
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

// validRefreshTokenRequest returns a RefreshTokenPayload with a sample token.
func validRefreshTokenRequest() *userauthpb.RefreshTokenPayload {
	return &userauthpb.RefreshTokenPayload{
		RefreshToken: "sample-refresh-token",
	}
}

// TestRefreshToken_Success ensures that a valid refresh token is successfully
// exchanged for a new token pair.
func TestRefreshToken_Success(t *testing.T) {
	// Scenario: Successful refresh token rotation.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshTokenRequest()
	user := testdata.UserProfileResponse()
	newTokenPair := testdata.ValidTokenPair()

	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{})
	tokenRepo.On("GetRefreshToken", mock.Anything, mock.Anything).Return("1", nil)
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(user, nil)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.Anything).Return(nil)
	jwtMock.On("CreateTokenPair", mock.Anything).Return(newTokenPair, nil)
	tokenRepo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)
	mapper.On("ToAuthTokenResponse", newTokenPair).Return(&userauthpb.AuthTokenResponse{})

	_, err := svc.RefreshToken(context.Background(), req)
	assert.NoError(t, err)
}

// TestRefreshToken_TokenNotFound ensures that a non-existent refresh token
// results in a NotFound gRPC error.
func TestRefreshToken_TokenNotFound(t *testing.T) {
	// Scenario: The provided refresh token is not found in the repository.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshTokenRequest()

	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{})
	tokenRepo.On("GetRefreshToken", mock.Anything, mock.Anything).Return("", errs.ErrTokenNotFound)

	_, err := svc.RefreshToken(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

// TestRefreshToken_FetchUserFails ensures that a failure to fetch the user
// after validating the token results in an Internal gRPC error.
func TestRefreshToken_FetchUserFails(t *testing.T) {
	// Scenario: The user associated with the token cannot be fetched.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshTokenRequest()

	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{})
	tokenRepo.On("GetRefreshToken", mock.Anything, mock.Anything).Return("1", nil)
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(transport.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := svc.RefreshToken(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}

// TestRefreshToken_DeleteFails ensures that a failure to delete the old token
// results in an Internal gRPC error.
func TestRefreshToken_DeleteFails(t *testing.T) {
	// Scenario: The old refresh token cannot be deleted from the repository.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshTokenRequest()
	user := testdata.UserProfileResponse()

	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{})
	tokenRepo.On("GetRefreshToken", mock.Anything, mock.Anything).Return("1", nil)
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(user, nil)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.Anything).Return(errs.ErrDBFailure)

	_, err := svc.RefreshToken(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}

// TestRefreshToken_JWTGenFails ensures that a failure to generate a new token
// pair results in an Internal gRPC error.
func TestRefreshToken_JWTGenFails(t *testing.T) {
	// Scenario: A new token pair cannot be generated after validating the old token.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshTokenRequest()
	user := testdata.UserProfileResponse()

	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{})
	tokenRepo.On("GetRefreshToken", mock.Anything, mock.Anything).Return("1", nil)
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(user, nil)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.Anything).Return(nil)
	jwtMock.On("CreateTokenPair", mock.Anything).Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := svc.RefreshToken(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}
