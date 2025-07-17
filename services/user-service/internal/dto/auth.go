package dto

// DTO used for data transfer over the HTTP

// RegisterRequest is what your HTTP handler binds on POST v1/auth/register
type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

// LoginRequest is what your HTTP handler binds on POST v1/auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LogoutResponse can be a simple confirmation
type LogoutResponse struct {
	Message string `json:"message"`
}
