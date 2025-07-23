// Package service provides the gRPC service implementations for chat rooms.
package service

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	chatpb "github.com/mamataliev-dev/social-platform/api/gen/chat/v1"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/mapper"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
)

// RoomService implements the ChatServiceServer for managing chat rooms.
type RoomService struct {
	chatpb.UnimplementedChatServiceServer

	roomRepo model.RoomRepository
	mapper   mapper.RoomMapper
	logger   *slog.Logger
}

// NewRoomService constructs a RoomService with the given dependencies.
//
//   - roomRepo: interface for persisting and retrieving rooms.
//   - mapper:   converts between gRPC messages and internal models.
//   - logger:   structured logger for diagnostics.
func NewRoomService(
	roomRepo model.RoomRepository,
	mapper mapper.RoomMapper,
	logger *slog.Logger,
) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
		mapper:   mapper,
		logger:   logger,
	}
}

// CreateRoom creates a new chat room based on the client request.
//
// It maps the incoming gRPC request to the internal Room model, calls the
// repository to persist the room, and returns a CreateRoomResponse message.
// If persistence fails, it logs the error with structured metadata and
// returns a gRPC Internal error status.
func (s *RoomService) CreateRoom(ctx context.Context, req *chatpb.CreateRoomRequest) (*chatpb.CreateRoomResponse, error) {
	roomModel := s.mapper.ToRoomModel(req)

	room, err := s.roomRepo.CreateRoom(ctx, roomModel)
	if err != nil {
		s.logger.Error("unable to create room", slog.Any("error", err))
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	resp := s.mapper.ToCreateRoomResponse(room)
	return resp, nil
}
