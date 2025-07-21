// Package mocks provides mock implementations of repository and service interfaces
// for unit testing the user-service.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
)

// TokenRepoMock is a mock implementation of the TokenRepository interface.
// It allows tests to simulate saving, retrieving, and deleting refresh tokens
// without a real database connection.
type TokenRepoMock struct {
	mock.Mock
}

// SaveRefreshToken simulates persisting a refresh token.
// Tests can configure this to return an error to test failure scenarios.
func (m *TokenRepoMock) SaveRefreshToken(ctx context.Context, in domain.SaveRefreshTokenInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

// GetRefreshToken simulates retrieving a refresh token by its value.
// It can be configured to return a user ID string or an error.
func (m *TokenRepoMock) GetRefreshToken(ctx context.Context, in transport.RefreshTokenRequest) (string, error) {
	args := m.Called(ctx, in)
	return args.String(0), args.Error(1)
}

// DeleteRefreshToken simulates deleting a refresh token.
// Tests can configure this to return an error to test failure scenarios.
func (m *TokenRepoMock) DeleteRefreshToken(ctx context.Context, in transport.RefreshTokenRequest) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}
