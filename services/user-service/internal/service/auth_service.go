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

type AuthService struct {
	userauthpb.UnimplementedAuthServiceServer
	authRepo  model.AuthRepository
	userRepo  model.UserRepository
	tokenRepo model.TokenRepository
	jwtGen    model.JWTGeneratorInterface
	hasher    security.Hasher
	converter mapper.Converter
}

func NewAuthService(
	authRepo model.AuthRepository,
	userRepo model.UserRepository,
	tokenRepo model.TokenRepository,
	jwtGen model.JWTGeneratorInterface,
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

var expiresAt = time.Now().Add(7 * 24 * time.Hour)

func (s *AuthService) Register(ctx context.Context, req *userauthpb.RegisterRequest) (*userauthpb.AuthTokenResponse, error) {
	hashedPwd, err := s.hasher.HashPassword(req.GetPassword())
	if err != nil {
		slog.Error("failed to hash password", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	mUser := s.converter.ToUserModel(req, hashedPwd)
	newUser, err := s.authRepo.CreateUser(ctx, mUser)
	if err != nil {
		slog.Error("failed to create new user", "err", err)
		if grpcCredentialsErr := checkUniqueCredentialsError(err); grpcCredentialsErr != nil {
			return nil, grpcCredentialsErr
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	slog.Info("registered new user", "username", newUser.Username, "id", newUser.ID)

	tokenPair, err := s.jwtGen.CreateTokenPair(createTokenRequest(newUser))
	if err != nil {
		slog.Error("failed to generate token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, saveTokenRequest(tokenPair, newUser, expiresAt))
	if err != nil {
		slog.Error("failed to save refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := s.converter.ToAuthTokenResponse(tokenPair)
	return resp, nil
}

func (s *AuthService) Login(ctx context.Context, req *userauthpb.LoginRequest) (*userauthpb.AuthTokenResponse, error) {
	mUser := s.converter.ToLoginRequest(req)

	user, err := s.authRepo.FetchUserByEmail(ctx, mUser.Email)
	if err != nil {
		slog.Error("failed to login", "err", err)
		return nil, utils.GrpcUserNotFoundError(err)
	}

	if err := s.hasher.VerifyPassword(user.PasswordHash, req.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, errs.ErrInvalidPassword.Error())
	}

	slog.Info("login user", "nickname", user.Nickname, "id", user.ID)

	tokenPair, err := s.jwtGen.CreateTokenPair(createTokenRequest(user))
	if err != nil {
		slog.Error("failed to generate token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, saveTokenRequest(tokenPair, user, expiresAt))
	if err != nil {
		slog.Error("failed to save refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := s.converter.ToAuthTokenResponse(tokenPair)
	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context, req *userauthpb.RefreshTokenPayload) (*userauthpb.LogoutResponse, error) {
	refreshToken := s.converter.ToGetRefreshTokenRequest(req)

	err := s.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, errs.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrTokenNotFound.Error())
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	msg := transport.LogoutResponse{
		Message: "Logout successfully",
	}

	return s.converter.ToLogoutResponse(msg), nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *userauthpb.RefreshTokenPayload) (*userauthpb.AuthTokenResponse, error) {
	input := s.converter.ToRefreshTokenRequest(req)

	userIDStr, err := s.tokenRepo.GetRefreshToken(ctx, input)
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

	reqID := transport.FetchUserByIDRequest{UserId: userID}
	userDTO, err := s.userRepo.FetchUserByID(ctx, reqID)
	if err != nil {
		slog.Error("failed to get user by ID", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := s.tokenRepo.DeleteRefreshToken(ctx, input); err != nil {
		slog.Error("failed to delete old refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	pairIn := domain.CreateTokenPairInput{
		UserID:   userDTO.ID,
		Nickname: userDTO.Nickname,
	}
	newPair, err := s.jwtGen.CreateTokenPair(pairIn)
	if err != nil {
		slog.Error("failed to generate new token pair", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	saveIn := domain.SaveRefreshTokenInput{
		UserID:    userDTO.ID,
		Token:     newPair.RefreshToken,
		ExpiresAt: expiresAt,
	}
	if err := s.tokenRepo.SaveRefreshToken(ctx, saveIn); err != nil {
		slog.Error("failed to save new refresh token", "err", err)
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := s.converter.ToAuthTokenResponse(newPair)
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

func createTokenRequest(user model.User) domain.CreateTokenPairInput {
	return domain.CreateTokenPairInput{
		UserID:   user.ID,
		Nickname: user.Nickname,
	}
}

func saveTokenRequest(token model.TokenPair, user model.User, expiresAt time.Time) domain.SaveRefreshTokenInput {
	return domain.SaveRefreshTokenInput{
		Token:     token.RefreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
}
