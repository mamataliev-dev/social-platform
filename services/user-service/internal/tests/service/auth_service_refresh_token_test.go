package service

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"testing"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
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
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshRequest()

	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{RefreshToken: req.RefreshToken})

	userID := int64(1)
	nickname := "tester"
	newPair := model.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	tokenRepo.On("GetRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return(strconv.FormatInt(userID, 10), nil)

	userRepo.On("FetchUserByID", mock.Anything, transport.FetchUserByIDRequest{UserId: userID}).
		Return(transport.UserProfileResponse{ID: userID, Nickname: nickname}, nil)

	tokenRepo.On("DeleteRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return(nil)

	jwtMock.On("CreateTokenPair", domain.CreateTokenPairInput{UserID: userID, Nickname: nickname}).
		Return(newPair, nil)

	tokenRepo.On("SaveRefreshToken", mock.Anything, mock.MatchedBy(func(in domain.SaveRefreshTokenInput) bool {
		return in.UserID == userID && in.Token == newPair.RefreshToken
	})).Return(nil)

	mapper.On("ToAuthTokenResponse", newPair).Return(
		&userauthpb.AuthTokenResponse{
			AccessToken:  newPair.AccessToken,
			RefreshToken: newPair.RefreshToken,
		},
	)

	resp, err := svc.RefreshToken(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, newPair.AccessToken, resp.AccessToken)
	assert.Equal(t, newPair.RefreshToken, resp.RefreshToken)

	mapper.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
}

func TestRefreshToken_TokenNotFound(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshRequest()
	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{RefreshToken: req.RefreshToken})

	tokenRepo.On("GetRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return("", errs.ErrTokenNotFound)

	_, err := svc.RefreshToken(context.Background(), req)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrTokenNotFound.Error(), st.Message())

	tokenRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRefreshToken_FetchUserFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshRequest()
	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{RefreshToken: req.RefreshToken})

	tokenRepo.On("GetRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return("1", nil)

	userRepo.On("FetchUserByID", mock.Anything, transport.FetchUserByIDRequest{UserId: 1}).
		Return(transport.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := svc.RefreshToken(context.Background(), req)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	mapper.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_DeleteFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshRequest()
	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{RefreshToken: req.RefreshToken})

	tokenRepo.On("GetRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return("1", nil)

	userRepo.On("FetchUserByID", mock.Anything, transport.FetchUserByIDRequest{UserId: 1}).
		Return(transport.UserProfileResponse{ID: 1, Nickname: "test"}, nil)

	tokenRepo.On("DeleteRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return(errs.ErrDBFailure)

	_, err := svc.RefreshToken(context.Background(), req)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	mapper.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_JWTGenFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRefreshRequest()
	mapper.On("ToRefreshTokenRequest", req).Return(transport.RefreshTokenRequest{RefreshToken: req.RefreshToken})

	tokenRepo.On("GetRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return("1", nil)

	userRepo.On("FetchUserByID", mock.Anything, transport.FetchUserByIDRequest{UserId: 1}).
		Return(transport.UserProfileResponse{ID: 1, Nickname: "test"}, nil)

	tokenRepo.On("DeleteRefreshToken", mock.Anything, transport.RefreshTokenRequest{RefreshToken: req.RefreshToken}).
		Return(nil)

	jwtMock.On("CreateTokenPair", domain.CreateTokenPairInput{UserID: 1, Nickname: "test"}).
		Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := svc.RefreshToken(context.Background(), req)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	mapper.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
}
