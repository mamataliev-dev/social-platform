// Package model defines the core domain types and repository interfaces
// for the chat service.
package model

import (
	"context"
	"time"
)

// Room represents a chat room between two participants.
type Room struct {
	ID            string
	InitiatorID   int64
	ParticipantID int64
	CreatedAt     time.Time
}

// RoomRepository defines persistence operations for Room.
// Implementations handle storage and retrieval of chat rooms.
type RoomRepository interface {
	// CreateRoom stores a new Room in the backing store and returns
	// the Room populated with ID and CreatedAt, or an error on failure.
	CreateRoom(ctx context.Context, room Room) (Room, error)
}
