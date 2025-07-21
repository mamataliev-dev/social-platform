// Package testdata provides reusable test data and helpers for user-service tests.
// It supports Single Responsibility and Open/Closed principles.
package testdata

import (
	"time"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

func ValidTokenPair() model.TokenPair {
	return model.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}
}

var TestPasswordHash = "$2a$10$test.hash.mock"

// SampleUserModel returns a fully populated
func SampleUserModel() model.User {
	return model.User{
		ID:           1,
		Username:     "Test User",
		Nickname:     "tester",
		PasswordHash: TestPasswordHash,
		Email:        "test@example.com",
		Bio:          "This is a test user",
		AvatarURL:    "https://test.com",
	}
}

func UserProfileResponse() transport.UserProfileResponse {
	return transport.UserProfileResponse{
		ID:        1,
		Username:  "Test User",
		Nickname:  "tester",
		Email:     "test@example.com",
		Bio:       "This is a test user",
		AvatarURL: "https://test.com",
	}
}
