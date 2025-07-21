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
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
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
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()
	user := testdata.SampleUserModel()

	hashedPassword := testdata.TestPasswordHash
	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	})

	tokenResp := testdata.ValidTokenPair()
	expectedRefreshToken := tokenResp.RefreshToken
	expectedExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user model.User) bool {
		return user.Email == req.Email &&
			user.Nickname == req.Nickname &&
			user.Username == req.Username &&
			user.AvatarURL == req.AvatarUrl &&
			user.Bio == req.Bio &&
			user.PasswordHash == testdata.TestPasswordHash
	})).Return(user, nil)

	jwtMock.On("CreateTokenPair", mock.MatchedBy(func(input domain.CreateTokenPairInput) bool {
		return input.UserID == user.ID &&
			input.Nickname == user.Nickname
	})).Return(tokenResp, nil)

	tokenRepo.On("SaveRefreshToken", mock.Anything, mock.MatchedBy(func(input domain.SaveRefreshTokenInput) bool {
		return input.UserID == user.ID &&
			input.Token == expectedRefreshToken &&
			input.ExpiresAt.Sub(expectedExpiresAt) < 2*time.Second
	})).Return(nil)

	mapper.On("ToAuthTokenResponse", tokenResp).Return(
		&userauthpb.AuthTokenResponse{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
		})

	resp, err := svc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokenResp.AccessToken, resp.AccessToken)
	assert.Equal(t, tokenResp.RefreshToken, resp.RefreshToken)

	authRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRegister_InternalDBError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)
	hashedPassword := testdata.TestPasswordHash

	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	})

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(input model.User) bool {
		return input.Email == req.Email &&
			input.Nickname == req.Nickname &&
			input.Username == req.Username &&
			input.AvatarURL == req.AvatarUrl &&
			input.Bio == req.Bio &&
			input.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrDBFailure)

	_, err := svc.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRegister_InternalError(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)
	hashedPassword := testdata.TestPasswordHash

	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	})

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user model.User) bool {
		return user.Email == req.Email &&
			user.Nickname == req.Nickname &&
			user.Username == req.Username &&
			user.AvatarURL == req.AvatarUrl &&
			user.Bio == req.Bio &&
			user.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrInternal)

	_, err := svc.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRegister_NicknameTaken(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)
	hashedPassword := testdata.TestPasswordHash

	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	})

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user model.User) bool {
		return user.Email == req.Email &&
			user.Nickname == req.Nickname &&
			user.Username == req.Username &&
			user.AvatarURL == req.AvatarUrl &&
			user.Bio == req.Bio &&
			user.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrNicknameTaken)

	_, err := svc.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, errs.ErrNicknameTaken.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRegister_EmailTaken(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)
	hashedPassword := testdata.TestPasswordHash

	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	})

	authRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user model.User) bool {
		return user.Email == req.Email &&
			user.Nickname == req.Nickname &&
			user.Username == req.Username &&
			user.AvatarURL == req.AvatarUrl &&
			user.Bio == req.Bio &&
			user.PasswordHash == testdata.TestPasswordHash
	})).Return(model.User{}, errs.ErrEmailTaken)

	_, err := svc.Register(context.Background(), req)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, errs.ErrEmailTaken.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRegister_HashPasswordFails(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()

	hasher.On("HashPassword", req.Password).Return("", errs.ErrHashingFailed)

	_, err := svc.Register(context.Background(), req)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestRegister_JWTGenerationFail(t *testing.T) {
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()
	user := testdata.SampleUserModel()

	hasher.On("HashPassword", req.Password).Return(testdata.TestPasswordHash, nil)

	hashedPassword := testdata.TestPasswordHash
	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	})

	authRepo.On("CreateUser", mock.Anything, mock.Anything).Return(user, nil)

	jwtMock.On("CreateTokenPair", mock.MatchedBy(func(input domain.CreateTokenPairInput) bool {
		return input.UserID == user.ID &&
			input.Nickname == user.Nickname
	})).Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := svc.Register(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	authRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
	mapper.AssertExpectations(t)
}
