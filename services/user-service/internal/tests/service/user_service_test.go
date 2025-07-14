package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/mamataliev-dev/social-platform/api/gen/user"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/service"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/mocks"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/tests/testdata"
)

func validFetchRequest() *userpb.FetchUserProfileByNicknameRequest {
	return &userpb.FetchUserProfileByNicknameRequest{
		Nickname: "test",
	}
}

func TestGetUserByNickname_Success(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	userService := service.NewUserService(repo)

	req := validFetchRequest()
	user := testdata.ValidUserDTO()

	repo.On("GetUserByNickname", mock.Anything, mock.MatchedBy(func(input model.GetUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(user, nil)

	resp, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, resp.User.Email)
	assert.Equal(t, user.Nickname, resp.User.Nickname)
	assert.Equal(t, user.UserName, resp.User.UserName)
	assert.Equal(t, user.Bio, resp.User.Bio)
	assert.Equal(t, user.AvatarURL, resp.User.AvatarUrl)

	repo.AssertExpectations(t)
}

func TestGetUserByNickname_InternalDBErr(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	userService := service.NewUserService(repo)

	req := validFetchRequest()

	repo.On("GetUserByNickname", mock.Anything, mock.MatchedBy(func(input model.GetUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(model.UserDTO{}, errs.ErrDBFailure)

	_, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	repo.AssertExpectations(t)
}

func TestGetUserByNickname_InternalErr(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	userService := service.NewUserService(repo)

	req := validFetchRequest()

	repo.On("GetUserByNickname", mock.Anything, mock.MatchedBy(func(input model.GetUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(model.UserDTO{}, errs.ErrInternal)

	_, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, errs.ErrInternal.Error(), st.Message())

	repo.AssertExpectations(t)
}

func TestGetUserByNickname_NotFound(t *testing.T) {
	repo := new(mocks.UserRepoMock)
	userService := service.NewUserService(repo)

	req := validFetchRequest()

	repo.On("GetUserByNickname", mock.Anything, mock.MatchedBy(func(input model.GetUserByNicknameInput) bool {
		return input.Nickname == req.Nickname
	})).Return(model.UserDTO{}, errs.ErrUserNotFound)

	_, err := userService.FetchUserProfileByNickname(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, errs.ErrUserNotFound.Error(), st.Message())

	repo.AssertExpectations(t)
}
