package mocks

import (
	"github.com/stretchr/testify/mock"

	chatpb "github.com/mamataliev-dev/social-platform/api/gen/chat/v1"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
)

// RoomMapperMock is a testify mock for the RoomMapper interface.
// It records calls and returns configured DTOs or models for mapping methods.
type RoomMapperMock struct {
	mock.Mock
}

// ToRoomDTO mocks the conversion from internal model to DTO.
func (m *RoomMapperMock) ToRoomDTO(r model.Room) dto.Room {
	args := m.Called(r)
	return args.Get(0).(dto.Room)
}

// FromRoomDTO mocks the conversion from DTO back to internal model.
func (m *RoomMapperMock) FromRoomDTO(d dto.Room) model.Room {
	args := m.Called(d)
	return args.Get(0).(model.Room)
}

// ToRoomModel mocks mapping a gRPC CreateRoomRequest into the internal Room model.
func (m *RoomMapperMock) ToRoomModel(req *chatpb.CreateRoomRequest) model.Room {
	args := m.Called(req)
	return args.Get(0).(model.Room)
}

// ToCreateRoomResponse mocks mapping an internal Room model to a gRPC CreateRoomResponse.
func (m *RoomMapperMock) ToCreateRoomResponse(room model.Room) *chatpb.CreateRoomResponse {
	args := m.Called(room)
	return args.Get(0).(*chatpb.CreateRoomResponse)
}
