package model

import (
	"context"
	"time"
)

// Model used for internal data transfer

type TokenPair struct {
	AccessToken  string    // JWT or opaque access token
	RefreshToken string    // opaque refresh token (UUID)
	ExpiresAt    time.Time // when the access token expires
}

type SaveRefreshTokenRequest struct {
	UserID    int64
	Token     string
	ExpiresAt time.Time
}

type RefreshTokenRequest struct {
	Token string
}

type CreateTokenPairRequest struct {
	UserID   int64
	Nickname string
}

// TokenRepository defines how we store and retrieve refresh tokens.
type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, input SaveRefreshTokenRequest) error
	GetRefreshToken(ctx context.Context, input RefreshTokenRequest) (userID int64, err error)
	DeleteRefreshToken(ctx context.Context, input RefreshTokenRequest) error
}

// JWTGenerator handles creation of token pairs.
type JWTGeneratorInterface interface {
	CreateTokenPair(input CreateTokenPairRequest) (TokenPair, error)
	GenerateRefreshToken() (string, error)
}
