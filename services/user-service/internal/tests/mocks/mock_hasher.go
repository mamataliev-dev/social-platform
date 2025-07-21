// Package mocks provides mock implementations of repository and service interfaces
// for unit testing the user-service.
package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockHasher is a mock implementation of the Hasher interface.
// It allows tests to simulate password hashing and verification without the
// computational overhead of real hashing algorithms.
type MockHasher struct {
	mock.Mock
}

// HashPassword simulates hashing a password string.
// Tests can configure this to return a specific hash or an error.
func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// VerifyPassword simulates comparing a plaintext password against a hash.
// It can be configured to return nil on success or an error on failure.
func (m *MockHasher) VerifyPassword(hash, password string) error {
	args := m.Called(hash, password)
	return args.Error(0)
}
