package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mamataliev-dev/social-platform/api/gen/user"
	userauthpb "github.com/mamataliev-dev/social-platform/api/gen/user_auth"
)

func MapRegisterRequestToDomainUser(req *userauthpb.RegisterRequest, hashedPassword string) UserDTO {
	return UserDTO{
		UserName:     req.GetUserName(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Nickname:     req.GetNickname(),
		Bio:          req.GetBio(),
		AvatarURL:    req.GetAvatarUrl(),
	}
}

func MapDomainUserToRegisterResponse(u UserDTO) *userauthpb.RegisterResponse {
	var createdAt *timestamppb.Timestamp

	if !u.CreatedAt.IsZero() {
		createdAt = timestamppb.New(u.CreatedAt)
	}

	return &userauthpb.RegisterResponse{
		User: &userpb.UserProfile{
			Id:        u.ID,
			UserName:  u.UserName,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Bio:       u.Bio,
			AvatarUrl: u.AvatarURL,
			CreatedAt: createdAt,
			UpdatedAt: nil,
			LastLogin: nil,
		},
	}
}

func MapLoginRequestToInput(req *userauthpb.LoginRequest) LoginInput {
	return LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func MapDomainUserToLoginResponse(u UserDTO) *userauthpb.LoginResponse {
	var createdAt, lastLogin, updatedAt *timestamppb.Timestamp

	if !u.CreatedAt.IsZero() {
		createdAt = timestamppb.New(u.CreatedAt)
	}
	if !u.LastLogin.IsZero() {
		lastLogin = timestamppb.New(u.LastLogin)
	}
	if !u.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(u.UpdatedAt)
	}

	return &userauthpb.LoginResponse{
		User: &userpb.UserProfile{
			Id:        u.ID,
			UserName:  u.UserName,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Bio:       u.Bio,
			AvatarUrl: u.AvatarURL,
			CreatedAt: createdAt,
			UpdatedAt: lastLogin,
			LastLogin: updatedAt,
		},
	}
}
