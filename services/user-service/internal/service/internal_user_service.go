// Package service implements the business logic for internal-only user profile
// retrieval. It follows SOLID principles by focusing on a single responsibility
// (fetching user data by ID) and relying on abstractions (repositories and mappers)
// to invert dependencies. This service is not exposed to the public gateway.
package service

import (
	"context"
	"log/slog"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

// InternalUserService handles internal, read-only access to user profiles by ID.
// It is used by other services within the system and depends on the UserRepository
// and Converter abstractions to remain decoupled from storage and transport details.
type InternalUserService struct {
	userpb.UnimplementedInternalUserServiceServer
	repo      model.UserRepository
	converter mapper.Converter
}

// NewInternalUserService constructs an InternalUserService with all required
// dependencies injected. This follows Dependency Inversion by relying on abstractions.
func NewInternalUserService(repo model.UserRepository, converter mapper.Converter) *InternalUserService {
	return &InternalUserService{
		repo:      repo,
		converter: converter,
	}
}

// FetchUserProfileByID retrieves a user profile by its unique numeric ID.
// It orchestrates data retrieval through the repository and maps the result to a
// gRPC response. Returns NotFound if the user does not exist or Internal on other
// failures.
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
