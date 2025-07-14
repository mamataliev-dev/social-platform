package model

import (
	"context"
	"time"
)

type UserDTO struct {
	ID           int64     `json:"id"`
	Nickname     string    `json:"nickname"`
	UserName     string    `json:"user_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserByNicknameInput struct {
	Nickname string `json:"nickname"`
}

type AuthRepository interface {
	Create(ctx context.Context, user UserDTO) (UserDTO, error)
	GetUserByEmail(ctx context.Context, input LoginInput) (UserDTO, error)
	GetUserByID(ctx context.Context, userID int64) (UserDTO, error)
}

type UserRepository interface {
	GetUserByNickname(ctx context.Context, input GetUserByNicknameInput) (UserDTO, error)
}
