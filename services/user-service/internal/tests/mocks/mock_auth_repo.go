package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type AuthRepoMock struct {
	mock.Mock
}

func (m *AuthRepoMock) Login(ctx context.Context, input transport.LoginRequest) (model.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *AuthRepoMock) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *AuthRepoMock) Logout(ctx context.Context, input transport.LoginRequest) (transport.LogoutResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(transport.LogoutResponse), args.Error(1)
}

func (m *AuthRepoMock) RefreshToken(ctx context.Context, input transport.RefreshTokenRequest) (transport.TokenResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(transport.TokenResponse), args.Error(1)
}

func (m *AuthRepoMock) FetchUserByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}
