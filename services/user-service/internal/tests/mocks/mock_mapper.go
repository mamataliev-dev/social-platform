// Package mocks provides mock implementations of repository and service interfaces
// for unit testing the user-service.
package mocks

import (
	"github.com/stretchr/testify/mock"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

// MockMapper is a mock implementation of the Converter interface.
// It allows tests to simulate data mapping between transport and domain models
// without depending on the concrete Mapper implementation.
type MockMapper struct {
	mock.Mock
}

// ToUserModel simulates mapping a RegisterRequest to a domain User model.
func (m *MockMapper) ToUserModel(req *userauthpb.RegisterRequest, pwd string) model.User {
	args := m.Called(req, pwd)
	return args.Get(0).(model.User)
}

// ToLoginRequest simulates mapping a LoginRequest to a transport LoginRequest DTO.
func (m *MockMapper) ToLoginRequest(req *userauthpb.LoginRequest) transport.LoginRequest {
	args := m.Called(req)
	return args.Get(0).(transport.LoginRequest)
}

// ToAuthTokenResponse simulates mapping a TokenPair to an AuthTokenResponse.
func (m *MockMapper) ToAuthTokenResponse(pair model.TokenPair) *userauthpb.AuthTokenResponse {
	args := m.Called(pair)
	return args.Get(0).(*userauthpb.AuthTokenResponse)
}

// ToFetchUserByNicknameRequest simulates mapping a gRPC request to a transport DTO.
func (m *MockMapper) ToFetchUserByNicknameRequest(req *userpb.FetchUserProfileByNicknameRequest) transport.FetchUserByNicknameRequest {
	args := m.Called(req)
	return args.Get(0).(transport.FetchUserByNicknameRequest)
}

// ToFetchUserByIDRequest simulates mapping a gRPC request to a transport DTO.
func (m *MockMapper) ToFetchUserByIDRequest(req *userpb.FetchUserProfileByIDRequest) transport.FetchUserByIDRequest {
	args := m.Called(req)
	return args.Get(0).(transport.FetchUserByIDRequest)
}

// ToFetchUserProfileResponse simulates mapping a transport DTO to a gRPC response.
func (m *MockMapper) ToFetchUserProfileResponse(u transport.UserProfileResponse) *userpb.UserProfile {
	args := m.Called(u)
	return args.Get(0).(*userpb.UserProfile)
}

// ToGetRefreshTokenRequest simulates mapping a gRPC request to a transport DTO.
func (m *MockMapper) ToGetRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	args := m.Called(req)
	return args.Get(0).(transport.RefreshTokenRequest)
}

// ToLogoutResponse simulates mapping a transport DTO to a gRPC response.
func (m *MockMapper) ToLogoutResponse(domain transport.LogoutResponse) *userauthpb.LogoutResponse {
	args := m.Called(domain)
	return args.Get(0).(*userauthpb.LogoutResponse)
}

// ToRefreshTokenRequest simulates mapping a gRPC request to a transport DTO.
func (m *MockMapper) ToRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	args := m.Called(req)
	return args.Get(0).(transport.RefreshTokenRequest)
}
