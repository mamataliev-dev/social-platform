package dto

import "time"

// Room represents the JSON payload for a chat room.
type Room struct {
	ID            string    `json:"id"`
	InitiatorID   int64     `json:"initiator_id"`
	ParticipantID int64     `json:"participant_id"`
	CreatedAt     time.Time `json:"created_at"`
}
