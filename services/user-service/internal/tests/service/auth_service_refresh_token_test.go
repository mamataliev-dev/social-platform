package service

import (
	"context"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"testing"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
)

func validRefreshRequest() *userauthpb.RefreshTokenPayload {
	return &userauthpb.RefreshTokenPayload{
		RefreshToken: "refresh-token-uuid",
	}
}

func TestRefreshToken_Success(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validRefreshRequest()
	userID := int64(1)
	nickname := "tester"
	newPair := model.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	tokenRepo.On("GetRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return(strconv.FormatInt(userID, 10), nil)
	authRepo.On("GetUserByID", mock.Anything, userID).Return(model.UserDTO{ID: userID, Nickname: nickname}, nil)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return(nil)
	jwtGen.On("CreateTokenPair", userID, nickname).Return(newPair, nil)
	tokenRepo.On("SaveRefreshToken", mock.Anything, userID, newPair.RefreshToken, mock.Anything).Return(nil)

	resp, err := authService.RefreshToken(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, newPair.AccessToken, resp.AccessToken)
	assert.Equal(t, newPair.RefreshToken, resp.RefreshToken)
}

func TestRefreshToken_DeleteFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validRefreshRequest()

	tokenRepo.On("GetRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return("1", nil)
	authRepo.On("GetUserByID", mock.Anything, int64(1)).Return(model.UserDTO{ID: 1, Nickname: "test"}, nil)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return(errs.ErrDBFailure)

	_, err := authService.RefreshToken(context.Background(), req)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}

func TestRefreshToken_TokenNotFound(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validRefreshRequest()
	tokenRepo.On("GetRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return("", errs.ErrTokenNotFound)

	_, err := authService.RefreshToken(context.Background(), req)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrTokenNotFound.Error(), st.Message())
}

func TestRefreshToken_JWTGenFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validRefreshRequest()
	tokenRepo.On("GetRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return("1", nil)
	authRepo.On("GetUserByID", mock.Anything, int64(1)).Return(model.UserDTO{ID: 1, Nickname: "test"}, nil)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, model.GetRefreshToken{
		RefreshToken: req.RefreshToken,
	}).Return(nil)
	jwtGen.On("CreateTokenPair", int64(1), "test").Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := authService.RefreshToken(context.Background(), req)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}
