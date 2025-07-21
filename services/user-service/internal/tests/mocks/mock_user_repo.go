// Package mocks provides mock implementations of repository and service interfaces
// for unit testing the user-service.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
)

// UserRepoMock is a mock implementation of the UserRepository interface.
// It allows tests to simulate fetching user profiles by nickname or ID without
// a real database connection.
type UserRepoMock struct {
	mock.Mock
}

// FetchUserByNickname simulates retrieving a user profile by their nickname.
// It can be configured to return a UserProfileResponse or an error.
func (m *UserRepoMock) FetchUserByNickname(
	ctx context.Context,
	in transport.FetchUserByNicknameRequest,
) (transport.UserProfileResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(transport.UserProfileResponse), args.Error(1)
}

// FetchUserByID simulates retrieving a user profile by their ID.
// It can be configured to return a UserProfileResponse or an error.
func (m *UserRepoMock) FetchUserByID(
	ctx context.Context,
	in transport.FetchUserByIDRequest,
) (transport.UserProfileResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(transport.UserProfileResponse), args.Error(1)
}
