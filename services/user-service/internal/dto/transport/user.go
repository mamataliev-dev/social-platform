package transport

import "time"

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

// FetchUserByNicknameRequest represents the HTTP path parameters for
// GET /v1/users/{nickname}, used to look up a public user profile.
type FetchUserByNicknameRequest struct {
	Nickname string `json:"nickname"`
}

// FetchUserByEmailRequest represents the HTTP path parameters for
// INTERNAL USE ONLY: called by AuthService.Login to load a user’s stored password hash.
type FetchUserByEmailRequest struct {
	Email string `json:"email"`
}

// FetchUserByIDRequest represents the HTTP path parameters for
// GET /v1/domain/users/{user_id}.
// Internal-only user profile lookups by ID
type FetchUserByIDRequest struct {
	UserId int64 `json:"user_id"`
}
