package testdata

import (
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
)

func ValidTokenPair() dto.TokenResponse {
	return dto.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}
}

var TestPasswordHash = "$2a$10$test.hash.mock"

func ValidUserProfileResponse() dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:        1,
		Username:  "Test User",
		Nickname:  "tester",
		Email:     "test@example.com",
		Bio:       "This is a test user",
		AvatarURL: "https://test.com",
	}
}
