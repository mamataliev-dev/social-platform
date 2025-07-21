package service

import (
	"context"
	"log/slog"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

type InternalUserService struct {
	userpb.UnimplementedInternalUserServiceServer
	repo      model.UserRepository
	converter mapper.Converter
}

func NewInternalUserService(repo model.UserRepository, converter mapper.Converter) *InternalUserService {
	return &InternalUserService{
		repo:      repo,
		converter: converter,
	}
}

func (s *InternalUserService) FetchUserProfileByID(ctx context.Context, req *userpb.FetchUserProfileByIDRequest) (*userpb.UserProfile, error) {
	mUser := s.converter.ToFetchUserByIDRequest(req)

	user, err := s.repo.FetchUserByID(ctx, mUser)
	if err != nil {
		slog.Error("failed to fetch user by id", "err", err)
		return nil, utils.GrpcUserNotFoundError(err)
	}

	slog.Info("fetched user by id", "nickname", user.Nickname, "user_id", user.ID)
	resp := s.converter.ToFetchUserProfileResponse(user)
	return resp, nil
}
