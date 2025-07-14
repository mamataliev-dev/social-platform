package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type AuthRepoMock struct {
	mock.Mock
}

func (m *AuthRepoMock) Create(ctx context.Context, user model.UserDTO) (model.UserDTO, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.UserDTO), args.Error(1)
}

func (m *AuthRepoMock) Login(ctx context.Context, input model.LoginInput) (model.UserDTO, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.UserDTO), args.Error(1)
}

func (m *AuthRepoMock) GetUserByEmail(ctx context.Context, input model.LoginInput) (model.UserDTO, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.UserDTO), args.Error(1)
}

func (m *AuthRepoMock) GetUserByID(ctx context.Context, userID int64) (model.UserDTO, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(model.UserDTO), args.Error(1)
}
