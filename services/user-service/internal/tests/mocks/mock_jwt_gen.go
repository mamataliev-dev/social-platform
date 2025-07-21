// Package mocks provides mock implementations of repository and service interfaces
// for unit testing the user-service.
package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

// JWTGeneratorMock is a mock implementation of the JWTGeneratorInterface.
// It allows tests to simulate JWT and refresh token generation without dealing
// with actual signing keys or random generation.
type JWTGeneratorMock struct {
	mock.Mock
}

// CreateTokenPair simulates the creation of a new access/refresh token pair.
// Tests can configure this to return a specific TokenPair or an error.
func (m *JWTGeneratorMock) CreateTokenPair(
	input domain.CreateTokenPairInput,
) (model.TokenPair, error) {
	args := m.Called(input)
	return args.Get(0).(model.TokenPair), args.Error(1)
}

// GenerateRefreshToken simulates the generation of a new refresh token.
func (m *JWTGeneratorMock) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
