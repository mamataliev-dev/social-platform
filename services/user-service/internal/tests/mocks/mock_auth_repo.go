package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type AuthRepoMock struct {
	mock.Mock
}

func (m *AuthRepoMock) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *AuthRepoMock) Login(ctx context.Context, input dto.LoginRequest) (model.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepoMock) FetchUserByEmail(ctx context.Context, input dto.FetchUserByEmailInput) (model.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.User), args.Error(1)
}
