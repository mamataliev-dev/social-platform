package service

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/security"
)

type AuthService struct {
	userauthpb.UnimplementedAuthServiceServer
	repo model.UserRepository
}

func NewAuthService(repo model.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req *userauthpb.RegisterRequest) (*userauthpb.RegisterResponse, error) {
	hashedPwd, err := security.HashPassword(req.GetPassword())
	if err != nil {
		slog.Error("failed to hash password", "err", err)
		return nil, status.Error(codes.Internal, "could not hash password")
	}

	mUser := model.MapRegisterRequestToDomainUser(req, hashedPwd)
	newUser, err := s.repo.Create(ctx, mUser)
	if err != nil {
		slog.Error("failed to create userpb", "err", err)
		if grpcErr := checkUniqueCredentialsError(err); grpcErr != nil {
			return nil, grpcErr
		}
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	slog.Info("registered new userpb", "username", newUser.UserName, "id", newUser.ID)

	respUser := model.MapDomainUserToRegisterResponse(newUser)
	return respUser, nil
}

func (s *AuthService) Login(ctx context.Context, req *userauthpb.LoginRequest) (*userauthpb.LoginResponse, error) {
	mUser := model.MapLoginRequestToInput(req)

	user, err := s.repo.Login(ctx, mUser)
	if err != nil {
		slog.Error("failed to login", "err", err)

		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, errs.ErrUserNotFound.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := security.VerifyPassword(user.PasswordHash, req.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, errs.ErrInvalidPassword.Error())
	}

	respUser := model.MapDomainUserToLoginResponse(user)
	return respUser, nil
}

func checkUniqueCredentialsError(err error) error {
	switch {
	case errors.Is(err, errs.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, errs.ErrNicknameTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	}
	return nil
}
