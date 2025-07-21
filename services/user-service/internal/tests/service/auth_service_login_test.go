// Package service_test verifies the behavior of AuthService’s login logic,
// covering success and failure scenarios for user authentication.
package service_test

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
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

// validLoginRequest returns a LoginRequest with well-formed, valid credentials.
func validLoginRequest() *userauthpb.LoginRequest {
	return &userauthpb.LoginRequest{
		Email:    "test@gmail.com",
		Password: "secure-password",
	}
}

// TestLogin_Success ensures that, given valid credentials, Login returns
// an AuthTokenResponse containing the expected access and refresh tokens.
func TestLogin_Success(t *testing.T) {
	// Scenario: Successful login with valid credentials.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLoginRequest()
	user := testdata.SampleUserModel()

	mapper.
		On("ToLoginRequest", req).
		Return(transport.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})

	tokenResp := testdata.ValidTokenPair()
	expectedRT := tokenResp.RefreshToken
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)

	authRepo.
		On("FetchUserByEmail", mock.Anything, req.Email).
		Return(user, nil)

	hasher.
		On("VerifyPassword", user.PasswordHash, req.Password).
		Return(nil)

	jwtMock.
		On("CreateTokenPair", mock.MatchedBy(func(input domain.CreateTokenPairInput) bool {
			return input.UserID == user.ID && input.Nickname == user.Nickname
		})).
		Return(tokenResp, nil)

	tokenRepo.
		On("SaveRefreshToken", mock.Anything, mock.MatchedBy(func(input domain.SaveRefreshTokenInput) bool {
			return input.UserID == user.ID &&
				input.Token == expectedRT &&
				input.ExpiresAt.Sub(expectedExpiry) < 2*time.Second
		})).
		Return(nil)

	mapper.
		On("ToAuthTokenResponse", tokenResp).
		Return(&userauthpb.AuthTokenResponse{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
		})

	resp, err := svc.Login(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, tokenResp.AccessToken, resp.AccessToken)
	assert.Equal(t, tokenResp.RefreshToken, resp.RefreshToken)
}

// TestLogin_InternalDBError ensures that a database failure during
// FetchUserByEmail is translated into an Internal gRPC error.
func TestLogin_InternalDBError(t *testing.T) {
	// Scenario: Login fails due to an internal database error.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLoginRequest()

	mapper.
		On("ToLoginRequest", req).
		Return(transport.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})

	authRepo.
		On("FetchUserByEmail", mock.Anything, req.Email).
		Return(model.User{}, errs.ErrDBFailure)

	_, err := svc.Login(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}

// TestLogin_InternalError ensures that a generic internal error during
// FetchUserByEmail is also translated into an Internal gRPC error.
func TestLogin_InternalError(t *testing.T) {
	// Scenario: Login fails due to an unexpected internal error.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLoginRequest()

	mapper.
		On("ToLoginRequest", req).
		Return(transport.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})

	authRepo.
		On("FetchUserByEmail", mock.Anything, req.Email).
		Return(model.User{}, errs.ErrInternal)

	_, err := svc.Login(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}

// TestLogin_UserNotFound ensures that a missing user leads to a NotFound
// gRPC error.
func TestLogin_UserNotFound(t *testing.T) {
	// Scenario: Login fails because the user with the given email is not found.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLoginRequest()

	mapper.
		On("ToLoginRequest", req).
		Return(transport.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})

	authRepo.
		On("FetchUserByEmail", mock.Anything, req.Email).
		Return(model.User{}, errs.ErrUserNotFound)

	_, err := svc.Login(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrUserNotFound.Error(), st.Message())
}

// TestLogin_InvalidPassword ensures that an incorrect password leads to
// an Unauthenticated gRPC error.
func TestLogin_InvalidPassword(t *testing.T) {
	// Scenario: Login fails because the provided password is incorrect.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLoginRequest()
	user := testdata.SampleUserModel()

	mapper.
		On("ToLoginRequest", req).
		Return(transport.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})

	authRepo.
		On("FetchUserByEmail", mock.Anything, req.Email).
		Return(user, nil)

	hasher.
		On("VerifyPassword", user.PasswordHash, req.Password).
		Return(errs.ErrInvalidPassword)

	_, err := svc.Login(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, errs.ErrInvalidPassword.Error(), st.Message())
}

// TestLogin_JWTGenerationFail ensures that a failure in token generation
// results in an Internal gRPC error.
func TestLogin_JWTGenerationFail(t *testing.T) {
	// Scenario: Login fails because JWT creation errors out.
	authRepo := new(mocks.AuthRepoMock)
	userRepo := new(mocks.UserRepoMock)
	tokenRepo := new(mocks.TokenRepoMock)
	jwtMock := new(mocks.JWTGeneratorMock)
	hasher := new(mocks.MockHasher)
	mapper := new(mocks.MockMapper)
	svc := service.NewAuthService(authRepo, userRepo, tokenRepo, jwtMock, hasher, mapper)

	req := validLoginRequest()
	user := testdata.SampleUserModel()

	mapper.
		On("ToLoginRequest", req).
		Return(transport.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})

	authRepo.
		On("FetchUserByEmail", mock.Anything, req.Email).
		Return(user, nil)

	hasher.
		On("VerifyPassword", user.PasswordHash, req.Password).
		Return(nil)

	jwtMock.
		On("CreateTokenPair", mock.MatchedBy(func(input domain.CreateTokenPairInput) bool {
			return input.UserID == user.ID && input.Nickname == user.Nickname
		})).
		Return(model.TokenPair{}, errs.ErrTokenSigningFailed)

	_, err := svc.Login(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())
}
