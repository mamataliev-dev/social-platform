package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/model"
)

// RoomRepoMock is a testify mock for the RoomRepository interface.
// It records calls and returns configured responses for CreateRoom.
type RoomRepoMock struct {
	mock.Mock
}

// CreateRoom mocks the repository method to create a chat room.
// It accepts a context and room model, then returns the mocked result and error.
func (m *RoomRepoMock) CreateRoom(ctx context.Context, room model.Room) (model.Room, error) {
	args := m.Called(ctx, room)
	return args.Get(0).(model.Room), args.Error(1)
}
