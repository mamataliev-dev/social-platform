package service

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatpb "github.com/mamataliev-dev/social-platform/api/gen/chat/v1"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/tests/mocks"
)

// validLoginRequest returns a well-formed CreateRoomRequest with preset IDs.
func validLoginRequest() *chatpb.CreateRoomRequest {
	return &chatpb.CreateRoomRequest{
		InitiatorId:   1,
		ParticipantId: 2,
	}
}

// validExpectedResp constructs the CreateRoomResponse matching validLoginRequest.
// It sets a fixed CreatedAt timestamp to compare against the service output.
func validExpectedResp() *chatpb.CreateRoomResponse {
	createdAt := time.Date(2025, 7, 23, 15, 0, 0, 0, time.UTC)

	return &chatpb.CreateRoomResponse{
		Room: &chatpb.Room{
			Id:            "room-uuid-123",
			InitiatorId:   1,
			ParticipantId: 2,
			CreatedAt:     timestamppb.New(createdAt),
		},
	}
}

// TestCreateRoom_Success verifies that CreateRoom returns a mapped response
// and no error when the repository successfully creates the room.
func TestCreateRoom_Success(t *testing.T) {
	roomRepo := new(mocks.RoomRepoMock)
	roomMapper := new(mocks.RoomMapperMock)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := service.NewRoomService(roomRepo, roomMapper, logger)

	req := validLoginRequest()
	expectedResp := validExpectedResp()

	inModel := model.Room{
		InitiatorID:   expectedResp.Room.InitiatorId,
		ParticipantID: expectedResp.Room.ParticipantId,
	}

	createdModel := model.Room{
		ID:            expectedResp.Room.Id,
		InitiatorID:   expectedResp.Room.InitiatorId,
		ParticipantID: expectedResp.Room.ParticipantId,
		CreatedAt:     expectedResp.Room.CreatedAt.AsTime(),
	}

	roomMapper.On("ToRoomModel", req).Return(inModel)

	roomRepo.On("CreateRoom", mock.Anything, inModel).Return(createdModel, nil)

	roomMapper.On("ToCreateRoomResponse", createdModel).Return(expectedResp)

	resp, err := svc.CreateRoom(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	roomMapper.AssertExpectations(t)
	roomRepo.AssertExpectations(t)
}

// TestCreateRoom_InternalError verifies that CreateRoom returns a gRPC Internal error
// when the repository fails to create the room.
func TestCreateRoom_InternalError(t *testing.T) {
	roomRepo := new(mocks.RoomRepoMock)
	roomMapper := new(mocks.RoomMapperMock)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := service.NewRoomService(roomRepo, roomMapper, logger)

	req := validLoginRequest()
	expectedResp := validExpectedResp()

	inModel := model.Room{
		InitiatorID:   expectedResp.Room.InitiatorId,
		ParticipantID: expectedResp.Room.ParticipantId,
	}

	createdModel := model.Room{
		ID:            expectedResp.Room.Id,
		InitiatorID:   expectedResp.Room.InitiatorId,
		ParticipantID: expectedResp.Room.ParticipantId,
		CreatedAt:     expectedResp.Room.CreatedAt.AsTime(),
	}

	roomMapper.On("ToRoomModel", req).Return(inModel)

	roomRepo.On("CreateRoom", mock.Anything, inModel).Return(createdModel, errs.ErrInternal)

	_, err := svc.CreateRoom(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	roomMapper.AssertExpectations(t)
	roomRepo.AssertExpectations(t)
}
