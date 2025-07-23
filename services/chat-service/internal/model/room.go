package model

import (
	"context"
	"time"
)

// Room represents a chat room between two participants.
// - ID: unique identifier assigned upon creation.
// - InitiatorID: user ID of the room creator.
// - ParticipantID: user ID of the other participant.
// - CreatedAt: timestamp when the room was created.
type Room struct {
	ID            string    // unique room UUID
	InitiatorID   int64     // creator's user ID
	ParticipantID int64     // other participant's user ID
	CreatedAt     time.Time // creation timestamp
}

// RoomRepository defines persistence operations for chat rooms.
// Implementers must handle storage and retrieval of Room entities.
type RoomRepository interface {
	// CreateRoom stores a new Room in the backing store and returns
	// the Room populated with ID and CreatedAt, or an error on failure.
	CreateRoom(ctx context.Context, room Room) (Room, error)
}
