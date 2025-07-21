// Package domain defines DTOs for domain-level token operations in the user-service.
// It supports Single Responsibility and Open/Closed principles.
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
