package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
)

type UserPostgres struct {
	DB *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{DB: db}
}

// FetchUserByNickname retrieves a user by their unique nickname.
func (r *UserPostgres) FetchUserByNickname(ctx context.Context, input dto.FetchUserByNicknameInput) (dto.UserProfileResponse, error) {
	const query = `
        SELECT id, user_name, email, nickname, bio, avatar_url, last_login_at, created_at, updated_at
        FROM users
        WHERE nickname = $1
    `
	row := r.DB.QueryRowContext(ctx, query, input.Nickname)

	return scanUserProfile(row)
}

// FetchUserByID retrieves a user by their id (internal uses only).
func (r *UserPostgres) FetchUserByID(ctx context.Context, input dto.FetchUserByIDInput) (dto.UserProfileResponse, error) {
	query := `
		SELECT id, user_name, email, nickname, bio, avatar_url, last_login_at, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	row := r.DB.QueryRowContext(ctx, query, input.UserId)

	return scanUserProfile(row)
}

// scanUserProfile scans a sql.Row into a UserProfileResponse.
func scanUserProfile(row *sql.Row) (dto.UserProfileResponse, error) {
	var u dto.UserProfileResponse
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Nickname,
		&u.Bio,
		&u.AvatarURL,
		&u.LastLogin,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return dto.UserProfileResponse{}, mapDBError(err)
	}
	return u, nil
}

func mapDBError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w", errs.ErrUserNotFound)
	}
	return errs.ErrDBFailure
}
