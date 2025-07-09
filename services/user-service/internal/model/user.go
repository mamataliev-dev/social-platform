package model

import (
	"context"
	"time"
)

type UserDTO struct {
	ID           int32
	Nickname     string
	UserName     string
	Email        string
	PasswordHash string
	Bio          string
	AvatarURL    string
	LastLogin    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RegisterInput struct {
	UserName  string
	Email     string
	Password  string
	Nickname  string
	Bio       string
	AvatarURL string
}

type LoginInput struct {
	Email    string
	Password string
}

type UserRepository interface {
	Create(ctx context.Context, user UserDTO) (UserDTO, error)
	Login(ctx context.Context, input LoginInput) (UserDTO, error)
}
