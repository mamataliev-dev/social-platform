package service

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

func validLogoutRequest() *userauthpb.RefreshTokenPayload {
	return &userauthpb.RefreshTokenPayload{
		RefreshToken: "1234-refresh-token-uuid",
	}
}

func TestLogout_Success(t *testing.T) {
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

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input transport.RefreshTokenRequest) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(nil)

	msg := transport.LogoutResponse{
		Message: "Logout successfully",
	}
	mapper.On("ToLogoutResponse", msg).
		Return(&userauthpb.LogoutResponse{
			Message: "Logout successfully",
		})

	resp, err := svc.Logout(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, msg.Message, resp.Message)

	tokenRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestLogout_InternalDBError(t *testing.T) {
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

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input transport.RefreshTokenRequest) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(errs.ErrDBFailure)

	_, err := svc.Logout(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	tokenRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestLogout_InternalError(t *testing.T) {
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

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input transport.RefreshTokenRequest) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(errs.ErrInternal)

	_, err := svc.Logout(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	tokenRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestLogout_TokenNotFound(t *testing.T) {
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

	tokenRepo.On("DeleteRefreshToken", mock.Anything, mock.MatchedBy(func(input transport.RefreshTokenRequest) bool {
		return input.RefreshToken == req.RefreshToken
	})).Return(errs.ErrTokenNotFound)

	_, err := svc.Logout(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrTokenNotFound.Error(), st.Message())

	tokenRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}
