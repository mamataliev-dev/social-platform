// Package mocks provides mock implementations of repository and service interfaces
// for unit testing the user-service. Each mock allows tests to configure expected
// inputs and outputs without relying on real dependencies.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

// AuthRepoMock is a mock implementation of the AuthRepository interface.
// It enables tests to simulate user‐creation and lookup behavior without touching
// the database or other external systems.
type AuthRepoMock struct {
	mock.Mock
}

// CreateUser simulates persisting a new user record.
// Tests should set up expected arguments and return values via On(...).Return(...).
func (m *AuthRepoMock) CreateUser(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

// FetchUserByEmail simulates retrieving a user by their email address.
// Tests should configure this mock to return either a valid User or an error.
func (m *AuthRepoMock) FetchUserByEmail(
	ctx context.Context,
	email string,
) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}
