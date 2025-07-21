package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user/v1"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/transport"
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
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(userRepo, mapper)

	req := validFetchUserByIDRequest()
	user := testdata.UserProfileResponse()

	// Expect conversion from protobuf request to domain DTO
	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{
		UserId: req.GetUserId(),
	})

	// Expect repository to return a valid user profile
	userRepo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input transport.FetchUserByIDRequest) bool {
		return input.UserId == req.UserId
	})).Return(user, nil)

	// Expect conversion from domain user DTO to protobuf response
	mapper.On("ToFetchUserProfileResponse", user).Return(&userpb.UserProfile{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarURL,
	})

	// Execute the service call
	resp, err := svc.FetchUserProfileByID(context.Background(), req)

	// Verify output
	assert.NoError(t, err)
	assert.Equal(t, user.ID, resp.UserId)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.Nickname, resp.Nickname)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, user.Bio, resp.Bio)
	assert.Equal(t, user.AvatarURL, resp.AvatarUrl)

	// Ensure all mocks were called as expected
	userRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

// TestFetchUserByID_InternalDBErr ensures that unexpected low-level database failures
// are translated into a gRPC domain error with a generic domain error message.
func TestFetchUserByID_InternalDBErr(t *testing.T) {
	UserRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(UserRepo, mapper)

	req := validFetchUserByIDRequest()

	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{
		UserId: req.GetUserId(),
	})

	UserRepo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input transport.FetchUserByIDRequest) bool {
		return input.UserId == req.UserId
	})).Return(transport.UserProfileResponse{}, errs.ErrDBFailure)

	_, err := svc.FetchUserProfileByID(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	UserRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

// TestFetchUserByID_InternalErr tests the case where the repository returns
// a generic domain service error (e.g., logic error or nil dereference).
// The service should return a gRPC Internal error code.
func TestFetchUserByID_InternalErr(t *testing.T) {
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(userRepo, mapper)

	req := validFetchUserByIDRequest()

	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{
		UserId: req.GetUserId(),
	})

	userRepo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input transport.FetchUserByIDRequest) bool {
		return input.UserId == req.UserId
	})).Return(transport.UserProfileResponse{}, errs.ErrInternal)

	_, err := svc.FetchUserProfileByID(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	userRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

// TestFetchUserByID_NotFound ensures that when a user is not found in the DB,
// the service returns a gRPC NotFound error code with the appropriate message.
func TestFetchUserByID_NotFound(t *testing.T) {
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(userRepo, mapper)

	req := validFetchUserByIDRequest()

	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{
		UserId: req.GetUserId(),
	})

	userRepo.On("FetchUserByID", mock.Anything, mock.MatchedBy(func(input transport.FetchUserByIDRequest) bool {
		return input.UserId == req.UserId
	})).Return(transport.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := svc.FetchUserProfileByID(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrUserNotFound.Error(), st.Message())

	userRepo.AssertExpectations(t)
	mapper.AssertExpectations(t)
}
