package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type TokenRepoMock struct {
	mock.Mock
}

func (m *TokenRepoMock) SaveRefreshToken(ctx context.Context, input model.SaveRefreshTokenRequest) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *TokenRepoMock) GetRefreshToken(ctx context.Context, input model.RefreshTokenRequest) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}

func (m *TokenRepoMock) DeleteRefreshToken(ctx context.Context, input model.RefreshTokenRequest) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
