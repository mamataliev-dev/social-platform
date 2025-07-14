package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/security"
)

type AuthService struct {
	userauthpb.UnimplementedAuthServiceServer
	authRepo  model.AuthRepository
	tokenRepo model.TokenRepository
	jwtGen    model.JWTGeneratorInterface
	hasher    security.Hasher
}

func NewAuthService(
	repo model.AuthRepository,
	tokenRepo model.TokenRepository,
	jwtGen model.JWTGeneratorInterface,
	hasher security.Hasher) *AuthService {
	return &AuthService{
		authRepo:  repo,
		tokenRepo: tokenRepo,
		jwtGen:    jwtGen,
		hasher:    hasher,
	}
}

var expiresAt = time.Now().Add(7 * 24 * time.Hour)

func (s *AuthService) Register(ctx context.Context, req *userauthpb.RegisterRequest) (*userauthpb.AuthTokenResponse, error) {
	hashedPwd, err := s.hasher.HashPassword(req.GetPassword())
	if err != nil {
		slog.Error("failed to hash password", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrHashingFailed.Error())
	}

	mUser := model.MapRegisterRequestToDomainUser(req, hashedPwd)
	newUser, err := s.authRepo.Create(ctx, mUser)
	if err != nil {
		slog.Error("failed to create user", "err", err)
		if grpcErr := checkUniqueCredentialsError(err); grpcErr != nil {
			return nil, grpcErr
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	slog.Info("registered new user", "username", newUser.UserName, "id", newUser.ID)

	tokenPair, err := s.jwtGen.CreateTokenPair(newUser.ID, newUser.Nickname)
	if err != nil {
		slog.Error("failed to generate token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, newUser.ID, tokenPair.RefreshToken, expiresAt)
	if err != nil {
		slog.Error("failed to save refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := model.MapRefreshTokenToAuthResponse(tokenPair.AccessToken, tokenPair.RefreshToken)
	return resp, nil
}

func (s *AuthService) Login(ctx context.Context, req *userauthpb.LoginRequest) (*userauthpb.AuthTokenResponse, error) {
	// Include this code when protoc-gen-validate is installed
	//if err := validation.ValidateNotEmpty(req.Email, req.Password); err != nil {
	//	return nil, status.Error(codes.InvalidArgument, err.Error())
	//}

	mUser := model.MapLoginRequestToInput(req)

	user, err := s.authRepo.GetUserByEmail(ctx, mUser)
	if err != nil {
		slog.Error("failed to login", "err", err)
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := s.hasher.VerifyPassword(user.PasswordHash, req.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, errs.ErrInvalidPassword.Error())
	}

	slog.Info("login user", "username", user.UserName, "id", user.ID)

	tokenPair, err := s.jwtGen.CreateTokenPair(user.ID, user.Nickname)
	if err != nil {
		slog.Error("failed to generate token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, user.ID, tokenPair.RefreshToken, expiresAt)
	if err != nil {
		slog.Error("failed to save refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := model.MapRefreshTokenToAuthResponse(tokenPair.AccessToken, tokenPair.RefreshToken)
	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context, req *userauthpb.RefreshTokenPayload) (*userauthpb.LogoutResponse, error) {
	refreshToken := model.MapRefreshTokenRequestToInput(req)

	err := s.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, errs.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrTokenNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := model.MapToLogoutResponse("logged out successfully")
	return resp, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *userauthpb.RefreshTokenPayload) (*userauthpb.AuthTokenResponse, error) {
	/*
		The user should be automatically redirected to the /auth/refresh page on the user's side (frontend)
		In case if user's token has expired
	*/
	refreshToken := model.MapRefreshTokenRequestToInput(req)

	userIDStr, err := s.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		slog.Error("failed to fetch refresh token", "err", err)
		if errors.Is(err, errs.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrTokenNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		slog.Error("invalid userID in refresh token", "userID", userIDStr, "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get user by ID", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := s.tokenRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		slog.Error("failed to delete old refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	tokenPair, err := s.jwtGen.CreateTokenPair(user.ID, user.Nickname)
	if err != nil {
		slog.Error("failed to generate new token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.tokenRepo.SaveRefreshToken(ctx, user.ID, tokenPair.RefreshToken, expiresAt); err != nil {
		slog.Error("failed to save new refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := model.MapRefreshTokenToAuthResponse(tokenPair.AccessToken, tokenPair.RefreshToken)
	return resp, nil
}

func checkUniqueCredentialsError(err error) error {
	switch {
	case errors.Is(err, errs.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, errs.ErrEmailTaken.Error())
	case errors.Is(err, errs.ErrNicknameTaken):
		return status.Error(codes.AlreadyExists, errs.ErrNicknameTaken.Error())
	}
	return nil
}
