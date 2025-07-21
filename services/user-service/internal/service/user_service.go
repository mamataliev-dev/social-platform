// Package service implements the business logic for public user profile retrieval.
// It follows SOLID principles by focusing on a single responsibility (fetching
// public user data) and relying on abstractions (repositories and mappers)
// to invert dependencies.
package service

import (
	"context"
	"log/slog"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

// UserService handles public, read-only access to user profiles by nickname.
// It depends on the UserRepository and Converter abstractions to keep its
// business logic decoupled from data storage and transport details.
type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo      model.UserRepository
	converter mapper.Converter
}

// NewUserService constructs a UserService with all required dependencies injected.
// This follows Dependency Inversion—high‐level logic depends on abstractions,
// not concrete implementations.
func NewUserService(repo model.UserRepository, converter mapper.Converter) *UserService {
	return &UserService{
		repo:      repo,
		converter: converter,
	}
}

// FetchUserProfileByNickname retrieves a public user profile by its unique nickname.
// It orchestrates data retrieval through the repository and maps the result to a
// gRPC response. Returns NotFound if the user does not exist or Internal on other
// failures.
func (s *UserService) FetchUserProfileByNickname(ctx context.Context, req *userpb.FetchUserProfileByNicknameRequest) (*userpb.UserProfile, error) {
	mUser := s.converter.ToFetchUserByNicknameRequest(req)

	user, err := s.repo.FetchUserByNickname(ctx, mUser)
	if err != nil {
		slog.Error("failed to fetch user byu nickname", "err", err)
		return nil, utils.GrpcUserNotFoundError(err)
	}

	slog.Info("fetched user by nickname", "username", user.Username, "nickname", user.Nickname)
	resp := s.converter.ToFetchUserProfileResponse(user)
	return resp, nil
}
