package model

import (
	"context"
	"time"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
)

// Model used for domain data transfer

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// TokenRepository defines how we store and retrieve refresh tokens.
type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, input domain.SaveRefreshTokenInput) error
	GetRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) (string, error)
	DeleteRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) error
}

// JWTGeneratorInterface handles creation of token pairs.
type JWTGeneratorInterface interface {
	CreateTokenPair(input domain.CreateTokenPairInput) (TokenPair, error)
	GenerateRefreshToken() (string, error)
}
