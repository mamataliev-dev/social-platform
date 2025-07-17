package dto

import "time"

// DTO used for data transfer over the HTTP

// RefreshRequest is what you bind on POST v1/auth/refresh or v1/auth/logout
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// TokenResponse is your standard response for register/login/refresh
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
