package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) VerifyPassword(hash, password string) error {
	args := m.Called(hash, password)
	return args.Error(0)
}
