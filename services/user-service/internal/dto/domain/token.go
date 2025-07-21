package domain

import "time"

type CreateTokenPairInput struct {
	UserID   int64
	Nickname string
}

type SaveRefreshTokenInput struct {
	UserID    int64
	Token     string
	ExpiresAt time.Time
}
