package service

import (
	"context"
	"log/slog"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo      model.UserRepository
	converter mapper.Converter
}

func NewUserService(repo model.UserRepository, converter mapper.Converter) *UserService {
	return &UserService{
		repo:      repo,
		converter: converter,
	}
}

func (s *UserService) FetchUserProfileByID(ctx context.Context, req *userpb.FetchUserProfileByIDRequest) (*userpb.UserProfile, error) {
	mUser := s.converter.ToFetchUserByIDInput(req)

	user, err := s.repo.FetchUserByID(ctx, mUser)
	if err != nil {
		slog.Error("failed to fetch user by id", "err", err)
		return nil, utils.GrpcError(err)
	}

	slog.Info("fetched user by id", "username", user.Username, "user_id", user.ID)
	resp := s.converter.ToFetchUserProfileResponse(user)
	return resp, nil
}

func (s *UserService) FetchUserProfileByNickname(ctx context.Context, req *userpb.FetchUserProfileByNicknameRequest) (*userpb.UserProfile, error) {
	mUser := s.converter.ToFetchUserByNicknameInput(req)

	user, err := s.repo.FetchUserByNickname(ctx, mUser)
	if err != nil {
		slog.Error("failed to fetch user byu nickname", "err", err)
		return nil, utils.GrpcError(err)
	}

	slog.Info("fetched user by nickname", "username", user.Username, "nickname", user.Nickname)
	resp := s.converter.ToFetchUserProfileResponse(user)
	return resp, nil
}
