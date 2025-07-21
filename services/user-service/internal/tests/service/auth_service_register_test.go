// Package service_test verifies the behavior of AuthService’s registration logic,
// covering success and failure scenarios.
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

// validRegisterRequest returns a RegisterRequest with well-formed, valid user data.
func validRegisterRequest() *userauthpb.RegisterRequest {
	return &userauthpb.RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Nickname: "newbie",
	}
}

// TestRegister_Success ensures that a valid registration request creates a user,
// generates tokens, and returns them in an AuthTokenResponse.
func TestRegister_Success(t *testing.T) {
	// Scenario: Successful user registration with valid data.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()
	hashedPassword := "hashed-password"
	createdUser := model.User{ID: 1, Nickname: req.Nickname}
	tokenPair := testdata.ValidTokenPair()

	hasher.On("HashPassword", req.Password).Return(hashedPassword, nil)
	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{})
	authRepo.On("CreateUser", mock.Anything, mock.Anything).Return(createdUser, nil)
	jwtMock.On("CreateTokenPair", mock.Anything).Return(tokenPair, nil)
	tokenRepo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)
	mapper.On("ToAuthTokenResponse", tokenPair).Return(&userauthpb.AuthTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})

	resp, err := svc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokenPair.AccessToken, resp.AccessToken)
	authRepo.AssertExpectations(t)
}

// TestRegister_HashingFails ensures that a password hashing failure
// results in an Internal gRPC error.
func TestRegister_HashingFails(t *testing.T) {
	// Scenario: Registration fails because password hashing returns an error.
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
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}

// TestRegister_CreateUserFails_EmailTaken ensures that an email taken error
// results in an AlreadyExists gRPC error.
func TestRegister_CreateUserFails_EmailTaken(t *testing.T) {
	// Scenario: Registration fails because the email is already in use.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()
	hashedPassword := "hashed-password"

	hasher.On("HashPassword", req.Password).Return(hashedPassword, nil)
	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{})
	authRepo.On("CreateUser", mock.Anything, mock.Anything).Return(model.User{}, errs.ErrEmailTaken)

	_, err := svc.Register(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, errs.ErrEmailTaken.Error(), st.Message())
}

// TestRegister_TokenGenerationFails ensures that a JWT creation failure
// results in an Internal gRPC error.
func TestRegister_TokenGenerationFails(t *testing.T) {
	// Scenario: Registration fails because token generation returns an error.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validRegisterRequest()
	hashedPassword := "hashed-password"
	createdUser := model.User{ID: 1, Nickname: req.Nickname}

	hasher.On("HashPassword", req.Password).Return(hashedPassword, nil)
	mapper.On("ToUserModel", req, hashedPassword).Return(model.User{})
	authRepo.On("CreateUser", mock.Anything, mock.Anything).Return(createdUser, nil)
	jwtMock.On("CreateTokenPair", mock.Anything).Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := svc.Register(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}
