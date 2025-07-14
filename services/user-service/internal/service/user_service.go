package service

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FetchUserProfileByNickname(ctx context.Context, req *userpb.FetchUserProfileByNicknameRequest) (*userpb.FetchUserProfileByNicknameResponse, error) {
	slog.Info("fetching user")

	mUser := model.MapFetchUserByNicknameRequestToInput(req)

	user, err := s.repo.GetUserByNickname(ctx, mUser)
	if err != nil {
		slog.Error("failed to fetch user", "err", err)

		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errs.ErrUserNotFound.Error())
		}

		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	slog.Info("fetched user", "username", user.UserName, "id", user.ID)
	resp := model.MapDomainUserToFetchUserByNicknameResponse(user)
	return resp, nil
}
