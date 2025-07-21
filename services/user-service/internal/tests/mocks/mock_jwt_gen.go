package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type JWTGeneratorMock struct {
	mock.Mock
}

func (m *JWTGeneratorMock) CreateTokenPair(input domain.CreateTokenPairInput) (model.TokenPair, error) {
	args := m.Called(input)
	return args.Get(0).(model.TokenPair), args.Error(1)
}

func (m *JWTGeneratorMock) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
