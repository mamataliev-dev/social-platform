// Package mapper provides abstractions and implementations for mapping between
// transport, domain, and persistence models. It supports Single Responsibility
// and Interface Segregation principles.
package mapper

import (
	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

// Converter defines methods for mapping between transport, domain, and persistence
// models. It enables Interface Segregation and Dependency Inversion for mapping logic.
type Converter interface {
	ToUserModel(*userauthpb.RegisterRequest, string) model.User
	ToLoginRequest(*userauthpb.LoginRequest) transport.LoginRequest
	ToAuthTokenResponse(model.TokenPair) *userauthpb.AuthTokenResponse
	ToFetchUserByNicknameRequest(*userpb.FetchUserProfileByNicknameRequest) transport.FetchUserByNicknameRequest
	ToFetchUserByIDRequest(*userpb.FetchUserProfileByIDRequest) transport.FetchUserByIDRequest
	ToFetchUserProfileResponse(transport.UserProfileResponse) *userpb.UserProfile
	ToGetRefreshTokenRequest(*userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest
	ToLogoutResponse(transport.LogoutResponse) *userauthpb.LogoutResponse
	ToRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest
}
