// Package service_test verifies the behavior of UserService’s logic for fetching
// users by nickname.
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

// validFetchUserByNicknameRequest returns a FetchUserProfileByNicknameRequest with a sample nickname.
func validFetchUserByNicknameRequest() *userpb.FetchUserProfileByNicknameRequest {
	return &userpb.FetchUserProfileByNicknameRequest{
		Nickname: "test-nickname",
	}
}

// TestFetchUserByNickname_Success ensures that a valid nickname returns the correct user profile.
func TestFetchUserByNickname_Success(t *testing.T) {
	// Scenario: A user profile is successfully fetched by its nickname.
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(userRepo, mapper)

	req := validFetchUserByNicknameRequest()
	user := testdata.UserProfileResponse()

	mapper.On("ToFetchUserByNicknameRequest", req).Return(transport.FetchUserByNicknameRequest{})
	userRepo.On("FetchUserByNickname", mock.Anything, mock.Anything).Return(user, nil)
	mapper.On("ToFetchUserProfileResponse", user).Return(&userpb.UserProfile{})

	_, err := svc.FetchUserProfileByNickname(context.Background(), req)
	assert.NoError(t, err)
}

// TestFetchUserByNickname_NotFound ensures that a non-existent nickname results
// in a NotFound gRPC error.
func TestFetchUserByNickname_NotFound(t *testing.T) {
	// Scenario: A user profile cannot be found for the given nickname.
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(userRepo, mapper)

	req := validFetchUserByNicknameRequest()

	mapper.On("ToFetchUserByNicknameRequest", req).Return(transport.FetchUserByNicknameRequest{})
	userRepo.On("FetchUserByNickname", mock.Anything, mock.Anything).Return(transport.UserProfileResponse{}, errs.ErrUserNotFound)

	_, err := svc.FetchUserProfileByNickname(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

// TestFetchUserByNickname_InternalError ensures that a generic repository error
// results in an Internal gRPC error.
func TestFetchUserByNickname_InternalError(t *testing.T) {
	// Scenario: An unexpected internal error occurs while fetching a user profile.
	userRepo := new(mocks.UserRepoMock)
	mapper := new(mocks.MockMapper)
	svc := service.NewUserService(userRepo, mapper)

	req := validFetchUserByNicknameRequest()

	mapper.On("ToFetchUserByNicknameRequest", req).Return(transport.FetchUserByNicknameRequest{})
	userRepo.On("FetchUserByNickname", mock.Anything, mock.Anything).Return(transport.UserProfileResponse{}, errs.ErrInternal)

	_, err := svc.FetchUserProfileByNickname(context.Background(), req)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}
