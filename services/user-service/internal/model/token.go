package model

import (
	"context"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token GetRefreshToken) (userID string, err error)
	DeleteRefreshToken(ctx context.Context, token GetRefreshToken) error
}

type JWTGeneratorInterface interface {
	CreateTokenPair(userID int64, nickname string) (TokenPair, error)
	GenerateRefreshToken() (string, error)
}
