package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

type AuthPostgres struct {
	DB *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{DB: db}
}

// CreateUser inserts a new user and handles unique constraint errors.
func (r *AuthPostgres) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	const query = `
		INSERT INTO users
		(username, email, password_hash, nickname, bio, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
    `

	err := r.DB.QueryRowContext(ctx, query,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.Nickname,
		u.Bio,
		u.AvatarURL,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				return model.User{}, fmt.Errorf("%w", errs.ErrEmailTaken)
			case "users_nickname_key":
				return model.User{}, fmt.Errorf("%w", errs.ErrNicknameTaken)
			default:
				return model.User{}, fmt.Errorf("duplicate constraint %q violated", pqErr.Constraint)
			}
		}

		return model.User{}, errs.ErrDBFailure
	}

	return u, nil
}

// FetchUserByEmail retrieves a user by their email address.
// INTERNAL USE ONLY: called by AuthService.Login to load a user’s stored password hash.
func (r *AuthPostgres) FetchUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT id, nickname, email, password_hash FROM users
		WHERE email = $1
	`

	var u model.User
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%w", errs.ErrUserNotFound)
		}
		return model.User{}, errs.ErrDBFailure
	}

	return u, nil
}
