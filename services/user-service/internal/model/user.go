package model

import (
	"context"
	"time"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
)

// ---------------------------------------------------------------------
// Model used for domain data transfer
// ---------------------------------------------------------------------

// User represents a domain user entity.
// Fields with pointer types are optional (maybe nil).
type User struct {
	ID           int64     // Unique identifier
	Username     string    // Chosen display name
	Email        string    // User's email address
	PasswordHash string    // Hashed password (never exposed)
	Nickname     string    // Unique nickname for public profile lookup
	Bio          string    // Public biography
	AvatarURL    string    // Avatar image URL
	LastLogin    time.Time // Timestamp of last login, if any
	CreatedAt    time.Time // Account creation timestamp
	UpdatedAt    time.Time // Timestamp of last profile update
}

// AuthRepository defines methods for user authentication persistence Create.
type AuthRepository interface {
	// CreateUser persists a new user and returns the created entity with ID and timestamps populated.
	CreateUser(ctx context.Context, user User) (User, error)

	// FetchUserByEmail retrieves a user by email; returns ErrUserNotFound if no record exists.
	// INTERNAL: used by AuthService.Login for password validation.
	FetchUserByEmail(ctx context.Context, email string) (User, error)
}

// UserRepository defines read-only retrieval by ID, Nickname or Email.
type UserRepository interface {
	// FetchUserByNickname looks up a public user profile by nickname.
	FetchUserByNickname(ctx context.Context, input transport.FetchUserByNicknameRequest) (transport.UserProfileResponse, error)

	// FetchUserByID retrieves a user by their numeric ID (domain uses only).
	FetchUserByID(ctx context.Context, input transport.FetchUserByIDRequest) (transport.UserProfileResponse, error)
}
