package mapper

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

// Mapper handles conversions between protobuf types and domain types.
type Mapper struct{}

// NewMapper creates a new Mapper instance.
func NewMapper() *Mapper {
	return &Mapper{}
}

// timestampOrNil converts a time.Time to a *timestamppb.Timestamp, returning nil if t is zero.
func timestampOrNil(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// ToUserModel maps a RegisterRequest and its hashed password to a User.
// Returns an empty User if req is nil.
func (m *Mapper) ToUserModel(req *userauthpb.RegisterRequest, hashedPassword string) model.User {
	if req == nil {
		return model.User{}
	}
	return model.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	}
}

// ToLoginRequest maps a LoginRequest to a LoginInput.
// Returns an empty LoginInput if req is nil.
func (m *Mapper) ToLoginRequest(req *userauthpb.LoginRequest) transport.LoginRequest {
	if req == nil {
		return transport.LoginRequest{}
	}
	return transport.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func (m *Mapper) ToAuthTokenResponse(pair model.TokenPair) *userauthpb.AuthTokenResponse {
	return &userauthpb.AuthTokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}

// ToFetchUserByNicknameRequest maps a FetchUserProfileByNicknameRequest to FetchUserByNicknameRequest.
// Returns an empty FetchUserByNicknameRequest if req is nil.
func (m *Mapper) ToFetchUserByNicknameRequest(req *userpb.FetchUserProfileByNicknameRequest) transport.FetchUserByNicknameRequest {
	if req == nil {
		return transport.FetchUserByNicknameRequest{}
	}
	return transport.FetchUserByNicknameRequest{Nickname: req.GetNickname()}
}

// ToFetchUserByIDRequest maps a FetchUserProfileByIDRequest to FetchUserByIDRequest.
// Returns an empty FetchUserByIDRequest if req is nil.
func (m *Mapper) ToFetchUserByIDRequest(req *userpb.FetchUserProfileByIDRequest) transport.FetchUserByIDRequest {
	if req == nil {
		return transport.FetchUserByIDRequest{}
	}
	return transport.FetchUserByIDRequest{UserId: req.GetUserId()}
}

// ToFetchUserProfileResponse maps a UserDTO to FetchUserProfileResponse.
func (m *Mapper) ToFetchUserProfileResponse(u transport.UserProfileResponse) *userpb.UserProfile {
	return &userpb.UserProfile{
		UserId:    u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Bio:       u.Bio,
		AvatarUrl: u.AvatarURL,
		CreatedAt: timestampOrNil(u.CreatedAt),
		UpdatedAt: timestampOrNil(u.UpdatedAt),
		LastLogin: timestampOrNil(u.LastLogin),
	}
}

// ToGetRefreshTokenRequest maps a RefreshTokenPayload to a GetRefreshToken.
// Returns an empty GetRefreshToken if req is nil.
func (m *Mapper) ToGetRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	if req == nil {
		return transport.RefreshTokenRequest{}
	}
	return transport.RefreshTokenRequest{RefreshToken: req.GetRefreshToken()}
}

// ToLogoutResponse creates a LogoutResponse with the provided message.
func (m *Mapper) ToLogoutResponse(domain transport.LogoutResponse) *userauthpb.LogoutResponse {
	return &userauthpb.LogoutResponse{
		Message: domain.Message,
	}
}

func (m *Mapper) ToRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	return transport.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}
