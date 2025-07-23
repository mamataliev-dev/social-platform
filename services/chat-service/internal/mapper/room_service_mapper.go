package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	chatpb "github.com/mamataliev-dev/social-platform/api/gen/chat/v1"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
)

// RoomMapper defines all mapping operations for Room.
type RoomMapper interface {
	// GRPC ↔ Domain
	ToRoomModel(req *chatpb.CreateRoomRequest) model.Room
	ToCreateRoomResponse(room model.Room) *chatpb.CreateRoomResponse

	// Domain ↔ DTO (JSON/DB)
	ToRoomDTO(room model.Room) dto.Room
	FromRoomDTO(d dto.Room) model.Room
}

type roomMapper struct{}

func NewRoomMapper() *roomMapper {
	return &roomMapper{}
}

// ToRoomDTO maps a domain model.Room into a persistence/JSON DTO.
func (m *roomMapper) ToRoomDTO(r model.Room) dto.Room {
	return dto.Room{
		ID:            r.ID,
		InitiatorID:   r.InitiatorID,
		ParticipantID: r.ParticipantID,
		CreatedAt:     r.CreatedAt,
	}
}

// FromRoomDTO maps a persistence/JSON DTO back into your domain model.Room.
func (m *roomMapper) FromRoomDTO(d dto.Room) model.Room {
	return model.Room{
		ID:            d.ID,
		InitiatorID:   d.InitiatorID,
		ParticipantID: d.ParticipantID,
		CreatedAt:     d.CreatedAt,
	}
}

// ToRoomModel (existing) maps the CreateRoomRequest into your domain model.Room.
func (m *roomMapper) ToRoomModel(req *chatpb.CreateRoomRequest) model.Room {
	if req == nil {
		return model.Room{}
	}
	return model.Room{
		InitiatorID:   req.GetInitiatorId(),
		ParticipantID: req.GetParticipantId(),
	}
}

// ToCreateRoomResponse (existing) maps your domain Room into the gRPC response.
func (m *roomMapper) ToCreateRoomResponse(r model.Room) *chatpb.CreateRoomResponse {
	return &chatpb.CreateRoomResponse{
		Room: &chatpb.Room{
			Id:            r.ID,
			InitiatorId:   r.InitiatorID,
			ParticipantId: r.ParticipantID,
			CreatedAt:     timestamppb.New(r.CreatedAt),
		},
	}
}
