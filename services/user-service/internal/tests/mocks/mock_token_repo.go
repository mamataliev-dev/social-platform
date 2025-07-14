package mocks

import (
	"context"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"time"

	"github.com/stretchr/testify/mock"
)

type TokenRepoMock struct {
	mock.Mock
}

func (m *TokenRepoMock) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *TokenRepoMock) GetRefreshToken(ctx context.Context, token model.GetRefreshToken) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *TokenRepoMock) DeleteRefreshToken(ctx context.Context, token model.GetRefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
