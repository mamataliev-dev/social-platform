package service

import (
	"context"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
)

var msg = "logged out successfully"

func validLogoutRequest() *userauthpb.RefreshTokenPayload {
	return &userauthpb.RefreshTokenPayload{
		RefreshToken: "1234-refresh-token-uuid",
	}
}

func TestLogout_Success(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher)

	req := validLogoutRequest()

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input model.GetRefreshToken) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(nil)

	resp, err := authService.Logout(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, msg, resp.Message)
	tokenRepo.AssertExpectations(t)
}

func TestLogout_InternalDBError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher)

	req := validLogoutRequest()

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input model.GetRefreshToken) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(errs.ErrDBFailure)

	_, err := authService.Logout(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}

func TestLogout_InternalError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher)

	req := validLogoutRequest()

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input model.GetRefreshToken) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(errs.ErrInternal)

	_, err := authService.Logout(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}

func TestLogout_TokenNotFound(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher)

	req := validLogoutRequest()

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input model.GetRefreshToken) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(errs.ErrTokenNotFound)

	_, err := authService.Logout(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrTokenNotFound.Error(), st.Message())
}

//ErrMissingOrInvalidAccessToken
//ErrMissingOrInvalidRefreshToken
