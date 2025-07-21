// Package model defines domain token entities and repository interfaces for
// refresh token and JWT management. It enables Dependency Inversion and Liskov
// Substitution for token storage and generation.
package model

import (
	"context"
	"time"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
)

// TokenPair represents a pair of access and refresh tokens with expiry.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// TokenRepository defines how we store and retrieve refresh tokens.
// It enables Dependency Inversion and Liskov Substitution for token storage.
type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, input domain.SaveRefreshTokenInput) error
	GetRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) (string, error)
	DeleteRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) error
}

// JWTGeneratorInterface handles creation of token pairs.
// It enables Dependency Inversion and Liskov Substitution for JWT generation.
type JWTGeneratorInterface interface {
	CreateTokenPair(input domain.CreateTokenPairInput) (TokenPair, error)
	GenerateRefreshToken() (string, error)
}
