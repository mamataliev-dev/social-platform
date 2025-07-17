package dto

import "time"

// DTO used for data transfer over the HTTP

// UserProfileResponse is returned on GET v1/users/{nickname}, v1/users/{email}, v1/users/{user_id}
type UserProfileResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}

// FetchUserByNicknameInput represents the HTTP path parameters for
// GET /v1/users/{nickname}, used to look up a public user profile.
type FetchUserByNicknameInput struct {
	Nickname string
}

// FetchUserByEmailInput represents the HTTP path parameters for
// INTERNAL USE ONLY: called by AuthService.Login to load a user’s stored password hash.
type FetchUserByEmailInput struct {
	Email string
}

// FetchUserByIDInput represents the HTTP path parameters for
// GET /v1/internal/users/{user_id}.
// Internal-only user profile lookups by ID
type FetchUserByIDInput struct {
	UserId int64
}
