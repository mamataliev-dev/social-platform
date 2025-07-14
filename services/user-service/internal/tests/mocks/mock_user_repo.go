package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) GetUserByNickname(ctx context.Context, input model.GetUserByNicknameInput) (model.UserDTO, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.UserDTO), args.Error(1)
}
