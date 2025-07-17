package repository

import (
	"context"
	"database/sql"

	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/utils"
)

type RoomPostgres struct {
	DB *sql.DB
}

func NewRoomPostgres(db *sql.DB) *RoomPostgres {
	return &RoomPostgres{DB: db}
}

func (r *RoomPostgres) CreatRoom(ctx context.Context, name string) error {}
