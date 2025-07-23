package dto

import "time"

// Room represents the JSON payload for a chat room.
// Fields are exported for JSON marshalling and annotated with
// `json` tags to define key names in HTTP responses.
type Room struct {
	ID            string    `json:"id"`             // Unique room identifier (UUID)
	InitiatorID   int64     `json:"initiator_id"`   // ID of the user who created the room
	ParticipantID int64     `json:"participant_id"` // ID of the other participant in the room
	CreatedAt     time.Time `json:"created_at"`     // Timestamp when the room was created
}
