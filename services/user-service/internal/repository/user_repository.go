package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/utils"
)

type UserPostgres struct {
	DB *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{DB: db}
}

func (r *UserPostgres) GetUserByNickname(ctx context.Context, input model.GetUserByNicknameInput) (model.UserDTO, error) {
	query := `
		SELECT id, user_name, email, nickname, bio, avatar_url, last_login_at, created_at, updated_at  FROM users 
		WHERE nickname = $1
	`

	var u model.UserDTO
	err := r.DB.QueryRow(query, input.Nickname).Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
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
