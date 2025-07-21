// Package transport defines DTOs for transport-level token operations in the
// user-service. It supports Single Responsibility and Open/Closed principles.
package transport

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
