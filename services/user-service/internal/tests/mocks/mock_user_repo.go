package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) FetchUserByNickname(ctx context.Context, input dto.FetchUserByNicknameInput) (dto.UserProfileResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(dto.UserProfileResponse), args.Error(1)
}

func (m *UserRepoMock) FetchUserByID(ctx context.Context, input dto.FetchUserByIDInput) (dto.UserProfileResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(dto.UserProfileResponse), args.Error(1)
}
