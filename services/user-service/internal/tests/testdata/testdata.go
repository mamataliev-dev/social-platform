package testdata

import (
	"time"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

func ValidTokenPair() model.TokenPair {
	return model.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}
}

var TestPasswordHash = "$2a$10$test.hash.mock"

func ValidUserDTO() model.UserDTO {
	now := time.Now()
	return model.UserDTO{
		ID:           1,
		Email:        "test@example.com",
		UserName:     "Test User",
		Nickname:     "tester",
		PasswordHash: TestPasswordHash,
		CreatedAt:    now,
	}
}
