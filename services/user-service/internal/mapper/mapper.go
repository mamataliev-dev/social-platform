package mapper

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
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

// ToUserDTO maps a RegisterRequest and its hashed password to a UserDTO.
// Returns an empty UserDTO if req is nil.
func (m *Mapper) ToUserDTO(req *userauthpb.RegisterRequest, hashedPassword string) model.User {
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

// ToLoginInput maps a LoginRequest to a LoginInput.
// Returns an empty LoginInput if req is nil.
func (m *Mapper) ToLoginInput(req *userauthpb.LoginRequest) dto.LoginRequest {
	if req == nil {
		return dto.LoginRequest{}
	}
	return dto.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

// ToAuthTokenResponse constructs an AuthTokenResponse from given access and refresh tokens.
func (m *Mapper) ToAuthTokenResponse(accessToken, refreshToken string) *userauthpb.AuthTokenResponse {
	return &userauthpb.AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// ToFetchUserByNicknameInput maps a FetchUserProfileByNicknameRequest to FetchUserByNicknameInput.
// Returns an empty FetchUserByNicknameInput if req is nil.
func (m *Mapper) ToFetchUserByNicknameInput(req *userpb.FetchUserProfileByNicknameRequest) dto.FetchUserByNicknameInput {
	if req == nil {
		return dto.FetchUserByNicknameInput{}
	}
	return dto.FetchUserByNicknameInput{Nickname: req.GetNickname()}
}

// ToFetchUserByIDInput maps a FetchUserProfileByIDRequest to FetchUserByIDInput.
// Returns an empty FetchUserByIDInput if req is nil.
func (m *Mapper) ToFetchUserByIDInput(req *userpb.FetchUserProfileByIDRequest) dto.FetchUserByIDInput {
	if req == nil {
		return dto.FetchUserByIDInput{}
	}
	return dto.FetchUserByIDInput{UserId: req.GetUserId()}
}

// ToFetchUserProfileResponse maps a UserDTO to FetchUserProfileResponse.
func (m *Mapper) ToFetchUserProfileResponse(u dto.UserProfileResponse) *userpb.UserProfile {
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

// ToGetRefreshToken maps a RefreshTokenPayload to a GetRefreshToken.
// Returns an empty GetRefreshToken if req is nil.
func (m *Mapper) ToGetRefreshToken(req *userauthpb.RefreshTokenPayload) dto.RefreshRequest {
	if req == nil {
		return dto.RefreshRequest{}
	}
	return dto.RefreshRequest{RefreshToken: req.GetRefreshToken()}
}

// ToLogoutResponse creates a LogoutResponse with the provided message.
func (m *Mapper) ToLogoutResponse(msg string) *userauthpb.LogoutResponse {
	return &userauthpb.LogoutResponse{Message: msg}
}
