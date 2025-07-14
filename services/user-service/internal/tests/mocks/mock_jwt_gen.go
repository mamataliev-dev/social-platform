package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type JWTGeneratorMock struct {
	mock.Mock
}

func (m *JWTGeneratorMock) CreateTokenPair(userID int64, nickname string) (model.TokenPair, error) {
	args := m.Called(userID, nickname)
	return args.Get(0).(model.TokenPair), args.Error(1)
}

func (m *JWTGeneratorMock) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
