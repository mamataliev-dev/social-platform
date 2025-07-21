// Package repository implements persistence logic for refresh token storage and
// retrieval. It provides concrete implementations of TokenRepository, following
// Dependency Inversion and Liskov Substitution principles.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
)

type TokenPostgres struct {
	DB *sql.DB
}

func NewTokenPostgres(db *sql.DB) *TokenPostgres {
	return &TokenPostgres{DB: db}
}

func (r *TokenPostgres) SaveRefreshToken(ctx context.Context, input domain.SaveRefreshTokenInput) error {
	query := `
        INSERT INTO refresh_tokens (token, user_id, expires_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (token) DO NOTHING
    `

	_, err := r.DB.ExecContext(ctx, query, input.Token, input.UserID, input.ExpiresAt)
	if err != nil {
		return errs.ErrDBFailure
	}

	return nil
}

func (r *TokenPostgres) GetRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) (string, error) {
	query := `
        SELECT user_id FROM refresh_tokens
        WHERE token = $1 AND expires_at > NOW()
    `

	var userID string
	err := r.DB.QueryRowContext(ctx, query, input.RefreshToken).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrTokenNotFound
		}
		return "", errs.ErrDBFailure
	}
	return userID, nil
}

func (r *TokenPostgres) DeleteRefreshToken(ctx context.Context, input transport.RefreshTokenRequest) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`

	result, err := r.DB.ExecContext(ctx, query, input.RefreshToken)
	if err != nil {
		return errs.ErrDBFailure
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errs.ErrTokenNotFound
	}

	return nil
}
