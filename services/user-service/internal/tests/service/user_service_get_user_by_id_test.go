// Package service_test verifies the behavior of InternalUserService’s logic for
// fetching users by ID.
package service_test

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

// validFetchUserByIDRequest returns a FetchUserProfileByIDRequest with a sample user ID.
func validFetchUserByIDRequest() *userpb.FetchUserProfileByIDRequest {
	return &userpb.FetchUserProfileByIDRequest{
		UserId: 1234,
	}
}

// TestFetchUserByID_Success ensures that a valid user ID returns the correct user profile.
func TestFetchUserByID_Success(t *testing.T) {
	// Scenario: A user profile is successfully fetched by its ID.
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewInternalUserService(userRepo, mapper)

	req := validFetchUserByIDRequest()
	user := testdata.UserProfileResponse()

	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{})
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(user, nil)
	mapper.On("ToFetchUserProfileResponse", user).Return(&userpb.UserProfile{})

	_, err := svc.FetchUserProfileByID(context.Background(), req)
	assert.NoError(t, err)
}

// TestFetchUserByID_NotFound ensures that a non-existent user ID results
// in a NotFound gRPC error.
func TestFetchUserByID_NotFound(t *testing.T) {
	// Scenario: A user profile cannot be found for the given ID.
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewInternalUserService(userRepo, mapper)

	req := validFetchUserByIDRequest()

	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{})
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(transport.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := svc.FetchUserProfileByID(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

// TestFetchUserByID_InternalError ensures that a generic repository error
// results in an Internal gRPC error.
func TestFetchUserByID_InternalError(t *testing.T) {
	// Scenario: An unexpected internal error occurs while fetching a user profile.
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewInternalUserService(userRepo, mapper)

	req := validFetchUserByIDRequest()

	mapper.On("ToFetchUserByIDRequest", req).Return(transport.FetchUserByIDRequest{})
	userRepo.On("FetchUserByID", mock.Anything, mock.Anything).Return(transport.UserProfileResponse{}, errs.ErrInternal)

	_, err := svc.FetchUserProfileByID(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}
