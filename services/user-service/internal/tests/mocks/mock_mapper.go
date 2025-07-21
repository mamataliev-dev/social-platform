package mocks

import (
	"github.com/stretchr/testify/mock"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type MockMapper struct {
	mock.Mock
}

func (m *MockMapper) ToUserModel(req *userauthpb.RegisterRequest, hashedPassword string) model.User {
	args := m.Called(req, hashedPassword)
	return args.Get(0).(model.User)
}

func (m *MockMapper) ToLoginRequest(req *userauthpb.LoginRequest) transport.LoginRequest {
	args := m.Called(req)
	return args.Get(0).(transport.LoginRequest)
}

func (m *MockMapper) ToAuthTokenResponse(pair model.TokenPair) *userauthpb.AuthTokenResponse {
	args := m.Called(pair)
	return args.Get(0).(*userauthpb.AuthTokenResponse)
}

func (m *MockMapper) ToFetchUserByNicknameRequest(req *userpb.FetchUserProfileByNicknameRequest) transport.FetchUserByNicknameRequest {
	args := m.Called(req)
	return args.Get(0).(transport.FetchUserByNicknameRequest)
}

func (m *MockMapper) ToFetchUserByIDRequest(req *userpb.FetchUserProfileByIDRequest) transport.FetchUserByIDRequest {
	args := m.Called(req)
	return args.Get(0).(transport.FetchUserByIDRequest)
}

func (m *MockMapper) ToFetchUserProfileResponse(u transport.UserProfileResponse) *userpb.UserProfile {
	args := m.Called(u)
	return args.Get(0).(*userpb.UserProfile)
}

func (m *MockMapper) ToGetRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	args := m.Called(req)
	return args.Get(0).(transport.RefreshTokenRequest)
}

func (m *MockMapper) ToLogoutResponse(msg transport.LogoutResponse) *userauthpb.LogoutResponse {
	args := m.Called(msg)
	return args.Get(0).(*userauthpb.LogoutResponse)
}

func (m *MockMapper) ToRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	args := m.Called(req)
	return args.Get(0).(transport.RefreshTokenRequest)
}

func (m *MockMapper) ToRefreshTokenResponse(domain transport.LogoutResponse) *userauthpb.LogoutResponse {
	args := m.Called(domain)
	return args.Get(0).(*userauthpb.LogoutResponse)
}
