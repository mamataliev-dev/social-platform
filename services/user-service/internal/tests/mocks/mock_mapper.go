package mocks

import (
	"github.com/stretchr/testify/mock"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type MockMapper struct {
	mock.Mock
}

func (m *MockMapper) ToUserDTO(req *userauthpb.RegisterRequest, hashedPassword string) model.User {
	args := m.Called(req, hashedPassword)
	return args.Get(0).(model.User)
}

func (m *MockMapper) ToLoginInput(req *userauthpb.LoginRequest) dto.LoginRequest {
	args := m.Called(req)
	return args.Get(0).(dto.LoginRequest)
}

func (m *MockMapper) ToAuthTokenResponse(accessToken, refreshToken string) *userauthpb.AuthTokenResponse {
	args := m.Called(accessToken, refreshToken)
	return args.Get(0).(*userauthpb.AuthTokenResponse)
}

func (m *MockMapper) ToFetchUserByNicknameInput(req *userpb.FetchUserProfileByNicknameRequest) dto.FetchUserByNicknameInput {
	args := m.Called(req)
	return args.Get(0).(dto.FetchUserByNicknameInput)
}

func (m *MockMapper) ToFetchUserByIDInput(req *userpb.FetchUserProfileByIDRequest) dto.FetchUserByIDInput {
	args := m.Called(req)
	return args.Get(0).(dto.FetchUserByIDInput)
}

func (m *MockMapper) ToFetchUserProfileResponse(u dto.UserProfileResponse) *userpb.UserProfile {
	args := m.Called(u)
	return args.Get(0).(*userpb.UserProfile)
}

func (m *MockMapper) ToGetRefreshToken(req *userauthpb.RefreshTokenPayload) dto.RefreshRequest {
	args := m.Called(req)
	return args.Get(0).(dto.RefreshRequest)
}

func (m *MockMapper) ToLogoutResponse(msg string) *userauthpb.LogoutResponse {
	args := m.Called(msg)
	return args.Get(0).(*userauthpb.LogoutResponse)
}
