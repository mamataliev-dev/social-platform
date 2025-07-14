package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user"
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

func MapLoginRequestToInput(req *userauthpb.LoginRequest) LoginInput {
	return LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func MapRefreshTokenToAuthResponse(accessToken, refreshToken string) *userauthpb.AuthTokenResponse {
	return &userauthpb.AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func MapFetchUserByNicknameRequestToInput(req *userpb.FetchUserProfileByNicknameRequest) GetUserByNicknameInput {
	return GetUserByNicknameInput{
		Nickname: req.GetNickname(),
	}
}

func MapDomainUserToFetchUserByNicknameResponse(u UserDTO) *userpb.FetchUserProfileByNicknameResponse {
	var createdAt, updatedAt, lastLogin *timestamppb.Timestamp

	if !u.CreatedAt.IsZero() {
		createdAt = timestamppb.New(u.CreatedAt)
	}
	if !u.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(u.UpdatedAt)
	}
	if !u.LastLogin.IsZero() {
		lastLogin = timestamppb.New(u.LastLogin)
	}

	return &userpb.FetchUserProfileByNicknameResponse{
		User: &userpb.UserProfile{
			Id:        u.ID,
			UserName:  u.UserName,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Bio:       u.Bio,
			AvatarUrl: u.AvatarURL,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			LastLogin: lastLogin,
		},
	}
}

func MapRefreshTokenRequestToInput(req *userauthpb.RefreshTokenPayload) GetRefreshToken {
	return GetRefreshToken{
		RefreshToken: req.GetRefreshToken(),
	}
}

func MapToLogoutResponse(msg string) *userauthpb.LogoutResponse {
	return &userauthpb.LogoutResponse{
		Message: msg,
	}
}
