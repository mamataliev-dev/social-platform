package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

func validRegisterRequest() *userauthpb.RegisterRequest {
	return &userauthpb.RegisterRequest{
		Email:     "test@example.com",
		Password:  "secure-password",
		Nickname:  "tester",
		Username:  "Test User",
		AvatarUrl: "https://test-avatar-url.com",
		Bio:       "test bio",
	}
}

func TestRegister_Success(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()
	user := testdata.ValidUserProfileResponse()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	tokenResp := testdata.ValidTokenPair()
	expectedRefreshToken := tokenResp.RefreshToken
	expectedExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(user, nil)

	jwtMock.On("CreateTokenPair", mock.Anything, mock.MatchedBy(func(input model.CreateTokenPairRequest) bool {
		return input.UserID == user.ID &&
			input.Nickname == user.Nickname
	})).Return(tokenResp, nil)

	tokenRepo.On("SaveRefreshToken", mock.Anything, mock.MatchedBy(func(input model.SaveRefreshTokenRequest) bool {
		return input.UserID == user.ID &&
			input.Token == expectedRefreshToken &&
			input.ExpiresAt.Sub(expectedExpiresAt) < 2*time.Second
	})).Return(nil)

	mapperMock.On("ToAuthTokenResponse", tokenResp.AccessToken, tokenResp.RefreshToken).Return(
		&userauthpb.AuthTokenResponse{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
		})

	resp, err := authService.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokenResp.AccessToken, resp.AccessToken)
	assert.Equal(t, tokenResp.RefreshToken, resp.RefreshToken)

	authRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestRegister_InternalDBError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrDBFailure)

	_, err := authService.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestRegister_InternalError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrInternal)

	_, err := authService.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestRegister_NicknameTaken(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrNicknameTaken)

	_, err := authService.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, errs.ErrNicknameTaken.Error(), st.Message())

	authRepo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestRegister_EmailTaken(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrEmailTaken)

	_, err := authService.Register(context.Background(), req)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, errs.ErrEmailTaken.Error(), st.Message())

	authRepo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestRegister_HashPasswordFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	hasher.On("HashPassword", req.Password).Return("", errs.ErrHashingFailed)

	_, err := authService.Register(context.Background(), req)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrHashingFailed.Error(), st.Message())

	hasher.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestRegister_JWTGenerationFail(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	tokenRepo := new(mocks.TokenRepoMock)
	hasher := new(mocks.MockHasher)
	mapperMock := new(mocks.MockMapper)
	authService := service.NewAuthService(authRepo, tokenRepo, jwtMock, hasher, mapperMock)

	req := validRegisterRequest()
	user := testdata.ValidUserProfileResponse()

	mapperMock.On("ToUserDTO", req).Return(dto.RegisterRequest{
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Nickname:  req.GetNickname(),
		Bio:       req.GetBio(),
		AvatarURL: req.GetAvatarUrl(),
	})

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, nil)

	jwtMock.On("CreateTokenPair", mock.Anything, mock.MatchedBy(func(input model.CreateTokenPairRequest) bool {
		return input.UserID == user.ID &&
			input.Nickname == user.Nickname
	})).Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := authService.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
}
