package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

func validLoginRequest() *userauthpb.LoginRequest {
	return &userauthpb.LoginRequest{
		Email:    "test@gmail.com",
		Password: "secure-password",
	}
}

func TestLogin_Success(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher)

	req := validLoginRequest()
	user := testdata.ValidUserProfileResponse()
	tokenResp := testdata.ValidTokenPair()

	expectedRefreshToken := tokenResp.RefreshToken
	expectedExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	authRepo.On("GetUserByEmail", mock.Anything, mock.MatchedBy(func(input model.LoginInput) bool {
		return input.Email == req.Email && input.Password == req.Password
	})).Return(user, nil)

	hasher.On("VerifyPassword", user.PasswordHash, req.Password).Return(nil)

	jwtMock.On("CreateTokenPair", user.ID, user.Nickname).Return(tokenResp, nil)

	tokenRepo.On("SaveRefreshToken", mock.Anything, user.ID, expectedRefreshToken, mock.MatchedBy(func(exp time.Time) bool {
		return exp.Sub(expectedExpiresAt) < 2*time.Second
	})).Return(nil)

	resp, err := authService.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokenResp.AccessToken, resp.AccessToken)
	assert.Equal(t, tokenResp.RefreshToken, resp.RefreshToken)

	authRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
}

func TestLogin_InternalDBError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validLoginRequest()

	authRepo.On("GetUserByEmail", mock.Anything, mock.MatchedBy(func(input model.LoginInput) bool {
		return input.Email == req.Email && input.Password == req.Password
	})).Return(model.UserDTO{}, errs.ErrDBFailure)

	_, err := authService.Login(context.Background(), req)

	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
}

func TestLogin_InternalError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validLoginRequest()

	authRepo.On("GetUserByEmail", mock.Anything, mock.MatchedBy(func(input model.LoginInput) bool {
		return input.Email == req.Email && input.Password == req.Password
	})).Return(model.UserDTO{}, errs.ErrInternal)

	_, err := authService.Login(context.Background(), req)

	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validLoginRequest()

	authRepo.On("GetUserByEmail", mock.Anything, mock.MatchedBy(func(input model.LoginInput) bool {
		return input.Email == req.Email && input.Password == req.Password
	})).Return(model.UserDTO{}, errs.ErrUserNotFound)

	_, err := authService.Login(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrUserNotFound.Error(), st.Message())

	authRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validLoginRequest()
	user := testdata.ValidUserProfileResponse()

	authRepo.On("GetUserByEmail", mock.Anything, mock.MatchedBy(func(input model.LoginInput) bool {
		return input.Email == req.Email && input.Password == req.Password
	})).Return(user, nil)

	hasher.On("VerifyPassword", user.PasswordHash, req.Password).Return(errs.ErrInvalidPassword)

	_, err := authService.Login(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, errs.ErrInvalidPassword.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

// Change this code when protoc-gen-third_party is installed
// Add input validation checker
//func TestLogin_InvalidInput(t *testing.T) {
//	authService := service.NewAuthService(nil, nil, nil)
//
//	req := &userauthpb.LoginRequest{
//		Email:    "",
//		Password: "",
//	}
//
//	_, err := authService.Login(context.Background(), req)
//
//	assert.Error(t, err)
//	st, ok := status.FromError(err)
//	assert.True(t, ok)
//	assert.Equal(t, codes.InvalidArgument, st.Code())
//	assert.Equal(t, errs.ErrMissingRequiredData.Error(), st.Message())
//}

func TestLogin_JWTGenerationFail(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtGen := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtGen, hasher)

	req := validLoginRequest()
	user := testdata.ValidUserProfileResponse()

	authRepo.On("GetUserByEmail", mock.Anything, mock.MatchedBy(func(input model.LoginInput) bool {
		return input.Email == req.Email && input.Password == req.Password
	})).Return(user, nil)

	hasher.On("VerifyPassword", user.PasswordHash, req.Password).Return(nil)

	jwtGen.On("CreateTokenPair", user.ID, user.Nickname).Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := authService.Login(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	jwtGen.AssertExpectations(t)
	hasher.AssertExpectations(t)
}
