package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
)

type TokenRepoMock struct {
	mock.Mock
}

func (m *TokenRepoMock) SaveRefreshToken(ctx context.Context, input domain.SaveRefreshTokenInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *TokenRepoMock) GetRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}

func (m *TokenRepoMock) DeleteRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
