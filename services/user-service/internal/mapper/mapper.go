// Package mapper provides a concrete implementation of the Converter interface,
// handling transformations between gRPC protobuf messages and internal domain models.
// It follows the Single Responsibility Principle by isolating data mapping logic,
// thus decoupling the service's transport layer from its business logic.
package mapper

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

// Mapper implements the Converter interface, providing methods to map data
// between gRPC/protobuf structures and the service's internal domain models.
// This decouples the transport layer from the business logic, adhering to the
// Single Responsibility and Dependency Inversion principles.
type Mapper struct{}

// NewMapper constructs a new Mapper instance. Its single responsibility is to
// provide a concrete implementation for data transformations.
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

// ToUserModel maps a gRPC RegisterRequest and a hashed password to a domain User model.
// This function's single responsibility is to handle the creation of a User
// entity from registration data, ensuring the transport model does not leak
// into the domain layer.
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

// ToLoginRequest maps a gRPC LoginRequest to a transport LoginRequest DTO.
// It isolates the transport-specific login structure from the service layer.
func (m *Mapper) ToLoginRequest(req *userauthpb.LoginRequest) transport.LoginRequest {
	if req == nil {
		return transport.LoginRequest{}
	}
	return transport.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

// ToAuthTokenResponse maps a domain TokenPair model to a gRPC AuthTokenResponse.
// This ensures the service's internal token representation is not directly
// exposed to the client.
func (m *Mapper) ToAuthTokenResponse(pair model.TokenPair) *userauthpb.AuthTokenResponse {
	return &userauthpb.AuthTokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}

// ToFetchUserByNicknameRequest maps a gRPC FetchUserProfileByNicknameRequest
// to a transport FetchUserByNicknameRequest DTO.
func (m *Mapper) ToFetchUserByNicknameRequest(req *userpb.FetchUserProfileByNicknameRequest) transport.FetchUserByNicknameRequest {
	if req == nil {
		return transport.FetchUserByNicknameRequest{}
	}
	return transport.FetchUserByNicknameRequest{Nickname: req.GetNickname()}
}

// ToFetchUserByIDRequest maps a gRPC FetchUserProfileByIDRequest to a
// transport FetchUserByIDRequest DTO.
func (m *Mapper) ToFetchUserByIDRequest(req *userpb.FetchUserProfileByIDRequest) transport.FetchUserByIDRequest {
	if req == nil {
		return transport.FetchUserByIDRequest{}
	}
	return transport.FetchUserByIDRequest{UserId: req.GetUserId()}
}

// ToFetchUserProfileResponse maps a transport UserProfileResponse DTO to a gRPC UserProfile.
// It handles the conversion of domain data, including timestamps, into the
// protobuf format expected by gRPC clients.
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

// ToGetRefreshTokenRequest maps a gRPC RefreshTokenPayload to a transport
// RefreshTokenRequest DTO.
func (m *Mapper) ToGetRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	if req == nil {
		return transport.RefreshTokenRequest{}
	}
	return transport.RefreshTokenRequest{RefreshToken: req.GetRefreshToken()}
}

// ToLogoutResponse maps a transport LogoutResponse DTO to a gRPC LogoutResponse.
func (m *Mapper) ToLogoutResponse(domain transport.LogoutResponse) *userauthpb.LogoutResponse {
	return &userauthpb.LogoutResponse{
		Message: domain.Message,
	}
}

// ToRefreshTokenRequest maps a gRPC RefreshTokenPayload to a transport
// RefreshTokenRequest DTO, used for token refresh operations.
func (m *Mapper) ToRefreshTokenRequest(req *userauthpb.RefreshTokenPayload) transport.RefreshTokenRequest {
	return transport.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}
