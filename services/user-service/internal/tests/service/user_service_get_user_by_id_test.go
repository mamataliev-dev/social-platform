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

// validFetchUserByIDRequest returns a sample request used in FetchUserProfileByID tests.
func validFetchUserByIDRequest() *userpb.FetchUserProfileByIDRequest {
	return &userpb.FetchUserProfileByIDRequest{
		UserId: 1234,
	}
}

// TestFetchUserByID_Success verifies that a valid user ID returns the correct user profile.
// It mocks both the mapper (for request/response transformation) and the repository.
func TestFetchUserByID_Success(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByIDRequest()
	user := testdata.ValidUserProfileResponse()

	// Expect conversion from protobuf request to internal DTO
	mapperMock.On("ToFetchUserByIDInput", req).Return(dto.FetchUserByIDInput{
		UserId: req.GetUserId(),
	})

	// Expect repository to return a valid user profile
	repo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByIDInput) bool {
		return input.UserId == req.UserId
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

	// Execute the service call
	resp, err := userService.FetchUserProfileByID(context.Background(), req)

	// Verify output
	assert.NoError(t, err)
	assert.Equal(t, user.ID, resp.UserId)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.Nickname, resp.Nickname)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, user.Bio, resp.Bio)
	assert.Equal(t, user.AvatarURL, resp.AvatarUrl)

	// Ensure all mocks were called as expected
	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

// TestFetchUserByID_InternalDBErr ensures that unexpected low-level database failures
// are translated into a gRPC internal error with a generic internal error message.
func TestFetchUserByID_InternalDBErr(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByIDRequest()

	mapperMock.On("ToFetchUserByIDInput", req).Return(dto.FetchUserByIDInput{
		UserId: req.GetUserId(),
	})

	repo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByIDInput) bool {
		return input.UserId == req.UserId
	})).Return(dto.UserProfileResponse{}, errs.ErrDBFailure)

	_, err := userService.FetchUserProfileByID(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

// TestFetchUserByID_InternalErr tests the case where the repository returns
// a generic internal service error (e.g., logic error or nil dereference).
// The service should return a gRPC Internal error code.
func TestFetchUserByID_InternalErr(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByIDRequest()

	mapperMock.On("ToFetchUserByIDInput", req).Return(dto.FetchUserByIDInput{
		UserId: req.GetUserId(),
	})

	repo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByIDInput) bool {
		return input.UserId == req.UserId
	})).Return(dto.UserProfileResponse{}, errs.ErrInternal)

	_, err := userService.FetchUserProfileByID(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

// TestFetchUserByID_NotFound ensures that when a user is not found in the DB,
// the service returns a gRPC NotFound error code with the appropriate message.
func TestFetchUserByID_NotFound(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	mapperMock := new(mocks.MockMapper)
	userService := service.NewUserService(repo, mapperMock)

	req := validFetchUserByIDRequest()

	mapperMock.On("ToFetchUserByIDInput", req).Return(dto.FetchUserByIDInput{
		UserId: req.GetUserId(),
	})

	repo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input dto.FetchUserByIDInput) bool {
		return input.UserId == req.UserId
	})).Return(dto.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := userService.FetchUserProfileByID(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrUserNotFound.Error(), st.Message())

	repo.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}
