package utils

import (
	"log/slog"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/clients"
)
    
func callUserClient() (*userpb.UserProfile, error) {
	// TODO: Add UserService addrs to the .env file, and use as a global variable
	uc, err := clients.NewUserServiceClient("5433:5432")
	if err != nil {
		slog.Error("failed to connect to user service", "error", err)
		return nil, err
	}

	user, err := uc.FetchUserByID(1)
	if err != nil {
		slog.Error("failed to fetch user", "error", err)
		return nil, err
	}

	return user, nil
}
