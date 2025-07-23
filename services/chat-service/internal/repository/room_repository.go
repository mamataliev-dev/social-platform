// Package repository provides PostgreSQL implementations of chat-service repositories.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
)

// RoomPostgres is a PostgreSQL implementation of model.RoomRepository.
type RoomPostgres struct {
	db     *sql.DB
	mapper mapper.RoomMapper
}

// NewRoomPostgres creates a new RoomPostgres backed by the given SQL DB.
func NewRoomPostgres(db *sql.DB, mapper mapper.RoomMapper) *RoomPostgres {
	return &RoomPostgres{
		db:     db,
		mapper: mapper,
	}
}

// CreateRoom inserts a new chat room record and returns the populated model.Room.
// It assigns a new UUID, persists the initiator and participant IDs, and
// populates CreatedAt. Returns ErrDBFailure on database errors.
func (r *RoomPostgres) CreateRoom(ctx context.Context, room model.Room) (model.Room, error) {
	id := uuid.New().String()
	room.ID = id

	dtoRoom := r.mapper.ToRoomDTO(room)

	query := `
        INSERT INTO rooms(id, initiator_id, participant_id)
        VALUES ($1, $2, $3)
        RETURNING id, initiator_id, participant_id, created_at
    `

	var inserted dto.Room
	err := r.db.QueryRowContext(ctx, query,
		dtoRoom.ID,
		dtoRoom.InitiatorID,
		dtoRoom.ParticipantID,
	).Scan(&inserted.ID, &inserted.InitiatorID, &inserted.ParticipantID, &inserted.CreatedAt)

	if err != nil {
		return model.Room{}, fmt.Errorf("%w: failed to insert room: %v", errs.ErrDBFailure, err)
	}

	return r.mapper.FromRoomDTO(inserted), nil
}
