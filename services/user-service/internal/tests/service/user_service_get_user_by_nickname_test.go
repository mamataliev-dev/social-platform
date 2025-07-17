package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

// validFetchUserByNicknameRequest returns a sample request used for fetching a user by nickname.
func validFetchUserByNicknameRequest() *userpb.FetchUserProfileByNicknameRequest {
	return &userpb.FetchUserProfileByNicknameRequest{
		Nickname: "test",
	}
}

// TestGetUserByNickname_Success verifies that a valid nickname returns the correct user profile.
// It mocks both the mapper (for input/output transformation) and the repository.
func TestGetUserByNickname_Success(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByNicknameRequest()
	user := testdata.ValidUserProfileResponse()

	// Expect conversion from protobuf request to internal DTO
	mapperMock.On("ToFetchUserByNicknameInput", req).Return(dto.FetchUserByNicknameInput{
		Nickname: req.GetNickname(),
	})

	// Expect repository to return a valid user profile
	repo.On("FetchUserByNickname", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(user, nil)

	// Expect conversion from internal user DTO to protobuf response
	mapperMock.On("ToFetchUserProfileResponse", user).Return(&userpb.UserProfile{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarURL,
	})

	// Execute the service method
	resp, err := userService.FetchUserProfileByNickname(context.Background(), req)

	// Validate the returned profile
	assert.NoError(t, err)
	assert.Equal(t, user.ID, resp.UserId)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.Nickname, resp.Nickname)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, user.Bio, resp.Bio)
	assert.Equal(t, user.AvatarURL, resp.AvatarUrl)

	// Ensure all mocks were triggered as expected
	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

// TestGetUserByNickname_InternalDBErr ensures that low-level DB errors
// are translated to gRPC internal errors with a generic internal message.
func TestGetUserByNickname_InternalDBErr(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByNicknameRequest()

	mapperMock.On("ToFetchUserByNicknameInput", req).Return(dto.FetchUserByNicknameInput{
		Nickname: req.GetNickname(),
	})

	repo.On("FetchUserByNickname", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(dto.UserProfileResponse{}, errs.ErrDBFailure)

	_, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

// TestGetUserByNickname_InternalErr tests handling of generic internal service errors.
// These should result in a gRPC Internal code and a general error message.
func TestGetUserByNickname_InternalErr(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByNicknameRequest()

	mapperMock.On("ToFetchUserByNicknameInput", req).Return(dto.FetchUserByNicknameInput{
		Nickname: req.GetNickname(),
	})

	repo.On("FetchUserByNickname", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(dto.UserProfileResponse{}, errs.ErrInternal)

	_, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

// TestGetUserByNickname_NotFound ensures that if the user is not found by nickname,
// the service responds with a gRPC NotFound error and the correct message.
func TestGetUserByNickname_NotFound(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByNicknameRequest()

	mapperMock.On("ToFetchUserByNicknameInput", req).Return(dto.FetchUserByNicknameInput{
		Nickname: req.GetNickname(),
	})

	repo.On("FetchUserByNickname", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(dto.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrUserNotFound.Error(), st.Message())

	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}
