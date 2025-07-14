package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

type AuthPostgres struct {
	DB *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{DB: db}
}

func (r *AuthPostgres) Create(ctx context.Context, u model.UserDTO) (model.UserDTO, error) {
	const query = `
		INSERT INTO users
		(user_name, email, password_hash, nickname, bio, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
    `

	err := r.DB.QueryRow(
		query,
		u.UserName,
		u.Email,
		u.PasswordHash,
		u.Nickname,
		u.Bio,
		u.AvatarURL,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		if ok, constraint := utils.IsUniqueViolation(err); ok {
			switch constraint {
			case "users_email_key":
				return model.UserDTO{}, fmt.Errorf("%w", errs.ErrEmailTaken)
			case "users_nickname_key":
				return model.UserDTO{}, fmt.Errorf("%w", errs.ErrNicknameTaken)
			default:
				return model.UserDTO{}, fmt.Errorf("duplicate constraint %q violated", constraint)
			}
		}
		return model.UserDTO{}, errs.ErrDBFailure
	}
	return u, nil
}

func (r *AuthPostgres) GetUserByEmail(ctx context.Context, input model.LoginInput) (model.UserDTO, error) {
	query := `
		SELECT id, user_name, email, password_hash, nickname, bio, avatar_url, last_login_at, created_at, updated_at FROM users 
		WHERE email = $1
	`

	var u model.UserDTO
	err := r.DB.QueryRow(query, input.Email).Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
		&u.PasswordHash,
		&u.Nickname,
		&u.Bio,
		&u.AvatarURL,
		&u.LastLogin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if !utils.IsUserExists(err) {
			return model.UserDTO{}, fmt.Errorf("%w", errs.ErrUserNotFound)
		}
		return model.UserDTO{}, errs.ErrDBFailure
	}

	return u, nil
}

func (r *AuthPostgres) GetUserByID(ctx context.Context, userId int64) (model.UserDTO, error) {
	query := `
		SELECT id, user_name, email, password_hash, nickname, bio, avatar_url, last_login_at, created_at, updated_at FROM users 
		WHERE id = $1
	`

	var u model.UserDTO
	err := r.DB.QueryRow(query, userId).Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
		&u.PasswordHash,
		&u.Nickname,
		&u.Bio,
		&u.AvatarURL,
		&u.LastLogin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if !utils.IsUserExists(err) {
			return model.UserDTO{}, fmt.Errorf("%w", errs.ErrUserNotFound)
		}
		return model.UserDTO{}, errs.ErrDBFailure
	}

	return u, nil
}
