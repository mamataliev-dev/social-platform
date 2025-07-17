package mapper

import (
	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type Converter interface {
	ToUserDTO(*userauthpb.RegisterRequest, string) model.User
	ToLoginInput(*userauthpb.LoginRequest) dto.LoginRequest
	ToAuthTokenResponse(string, string) *userauthpb.AuthTokenResponse
	ToFetchUserByNicknameInput(*userpb.FetchUserProfileByNicknameRequest) dto.FetchUserByNicknameInput
	ToFetchUserByIDInput(*userpb.FetchUserProfileByIDRequest) dto.FetchUserByIDInput
	ToFetchUserProfileResponse(dto.UserProfileResponse) *userpb.UserProfile
	ToGetRefreshToken(*userauthpb.RefreshTokenPayload) dto.RefreshRequest
	ToLogoutResponse(string) *userauthpb.LogoutResponse
}
