// Package service implements the business logic for the user-service.
// It contains services for authentication (AuthService), public user profile
// retrieval (UserService), and internal user lookups (InternalUserService).
// The package follows SOLID principles by keeping each operation focused on a single
// responsibility and relying on abstractions (repositories, JWT generators,
// hashers, and mappers) to invert dependencies.
package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/security"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

// AuthService orchestrates user registration, login, logout, and token refresh.
// It depends on AuthRepository, UserRepository, TokenRepository, JWTGenerator,
// Hasher, and Converter abstractions to keep business rules decoupled from storage.
type AuthService struct {
	authRepo  model.AuthRepository
	userRepo  model.UserRepository
	tokenRepo model.TokenRepository
	jwtGen    security.JWTGenerator
	hasher    security.Hasher
	converter mapper.Converter
}

// NewAuthService constructs an AuthService with all required dependencies injected.
// This follows Dependency Inversion—high‐level logic depends on abstractions.
func NewAuthService(
	authRepo model.AuthRepository,
	userRepo model.UserRepository,
	tokenRepo model.TokenRepository,
	jwtGen security.JWTGenerator,
	hasher security.Hasher,
	converter mapper.Converter) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtGen:    jwtGen,
		hasher:    hasher,
		converter: converter,
	}
}

// defaultTTL is the lifespan of a refresh token (7 days).
var defaultTTL = 7 * 24 * time.Hour

// Register creates a new user, generates an access/refresh token pair,
// and persists the refresh token. Errors during hashing, persistence,
// or token generation produce appropriate gRPC status codes.
func (s *AuthService) Register(
	ctx context.Context,
	req *userauthpb.RegisterRequest,
) (*userauthpb.AuthTokenResponse, error) {
	hashedPwd, err := s.hasher.HashPassword(req.GetPassword())
	if err != nil {
		slog.Error("failed to hash password", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	userModel := s.converter.ToUserModel(req, hashedPwd)
	createdUser, err := s.authRepo.CreateUser(ctx, userModel)
	if err != nil {
		slog.Error("failed to create new user", "err", err)
		if grpcErr := mapCreateUserError(err); grpcErr != nil {
			return nil, grpcErr
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	tokenPair, err := s.jwtGen.CreateTokenPair(domain.CreateTokenPairInput{
		UserID:   createdUser.ID,
		Nickname: createdUser.Nickname,
	})
	if err != nil {
		slog.Error("failed to generate token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := s.tokenRepo.SaveRefreshToken(ctx, domain.SaveRefreshTokenInput{
		UserID:    createdUser.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(defaultTTL),
	}); err != nil {
		slog.Error("failed to save refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	return s.converter.ToAuthTokenResponse(tokenPair), nil
}

// Login validates user credentials, issues a new token pair, and stores the
// refresh token. Returns Unauthenticated on bad password, NotFound if user
// doesn’t exist, or Internal on other failures.
func (s *AuthService) Login(
	ctx context.Context,
	req *userauthpb.LoginRequest,
) (*userauthpb.AuthTokenResponse, error) {
	loginInput := s.converter.ToLoginRequest(req)
	user, err := s.authRepo.FetchUserByEmail(ctx, loginInput.Email)
	if err != nil {
		slog.Error("login failed: user not found", "err", err)
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := s.hasher.VerifyPassword(user.PasswordHash, req.GetPassword()); err != nil {
		slog.Error("invalid password provided", "err", err)
		return nil, status.Error(codes.Unauthenticated, errs.ErrInvalidPassword.Error())
	}

	tokenPair, err := s.jwtGen.CreateTokenPair(domain.CreateTokenPairInput{
		UserID:   user.ID,
		Nickname: user.Nickname,
	})
	if err != nil {
		slog.Error("failed to generate token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := s.tokenRepo.SaveRefreshToken(ctx, domain.SaveRefreshTokenInput{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(defaultTTL),
	}); err != nil {
		slog.Error("failed to save refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	slog.Info("login successful", "email", req.GetEmail())
	return s.converter.ToAuthTokenResponse(tokenPair), nil
}

// Logout deletes an existing refresh token, returning NotFound if absent,
// or Internal on storage errors.
func (s *AuthService) Logout(
	ctx context.Context,
	req *userauthpb.RefreshTokenPayload,
) (*userauthpb.LogoutResponse, error) {
	input := s.converter.ToGetRefreshTokenRequest(req)
	if err := s.tokenRepo.DeleteRefreshToken(ctx, input); err != nil {
		slog.Error("failed to delete refresh token", "err", err)
		if errors.Is(err, errs.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrTokenNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	slog.Info("logout successful")
	return s.converter.ToLogoutResponse(transport.LogoutResponse{
		Message: "Logout successful",
	}), nil
}

// RefreshToken validates an old refresh token, loads the user profile,
// revokes the old token, issues a new token pair, and persists the new
// refresh token. It returns NotFound, Internal, or a valid AuthTokenResponse.
func (s *AuthService) RefreshToken(
	ctx context.Context,
	req *userauthpb.RefreshTokenPayload,
) (*userauthpb.AuthTokenResponse, error) {
	input := s.converter.ToRefreshTokenRequest(req)

	// 1) retrieve and validate existing token
	userIDStr, err := s.tokenRepo.GetRefreshToken(ctx, input)
	if err != nil {
		slog.Error("failed to fetch refresh token", "err", err)
		if errors.Is(err, errs.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrTokenNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	// 2) parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		slog.Error("invalid userID in refresh token", "userID", userIDStr, "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	// 3) load user profile
	profileReq := transport.FetchUserByIDRequest{UserId: userID}
	userDTO, err := s.userRepo.FetchUserByID(ctx, profileReq)
	if err != nil {
		slog.Error("failed to get user by ID", "err", err)
		return nil, utils.GrpcUserNotFoundError(err)
	}

	// 4) revoke old token
	if err := s.tokenRepo.DeleteRefreshToken(ctx, input); err != nil {
		slog.Error("failed to delete old refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	// 5) issue and persist new token pair
	newPair, err := s.jwtGen.CreateTokenPair(domain.CreateTokenPairInput{
		UserID:   userDTO.ID,
		Nickname: userDTO.Nickname,
	})
	if err != nil {
		slog.Error("failed to generate new token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}
	expiresAt := time.Now().Add(defaultTTL)
	saveIn := domain.SaveRefreshTokenInput{
		Token:     newPair.RefreshToken,
		UserID:    userDTO.ID,
		ExpiresAt: expiresAt,
	}

	if err = s.tokenRepo.SaveRefreshToken(ctx, saveIn); err != nil {
		slog.Error("failed to save new refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}
	slog.Info("token refreshed successfully", "userID", userDTO.ID)
	return s.converter.ToAuthTokenResponse(newPair), nil
}

// mapCreateUserError inspects repository errors for unique‐constraint violations
// and returns the appropriate gRPC AlreadyExists status. Returns nil otherwise.
func mapCreateUserError(err error) error {
	switch {
	case errors.Is(err, errs.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, errs.ErrEmailTaken.Error())
	case errors.Is(err, errs.ErrNicknameTaken):
		return status.Error(codes.AlreadyExists, errs.ErrNicknameTaken.Error())
	}
	return nil
}
