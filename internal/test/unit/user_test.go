package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/service"
	mockrepository "github.com/kidx45/Debter/internal/test/mock/repository"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func RandomUser(t *testing.T) domain.User {
	return domain.User{
		ID:             util.RandomNumber(1, 1000),
		FullName:       util.RandomFullName(6, 12),
		Username:       util.RandomUserName(6, 12),
		Email:          util.RandomEmail(6, 12),
		HashedPassword: util.RandomPassword(8, 16),
		CreatedAt:      time.Now(),
	}
}

func NewTestUserService(t *testing.T, userRepo *mockrepository.MockUserRepository) *service.UserService {
	return service.NewUserService(userRepo)
}

func TestCreateUser(t *testing.T) {
	user := RandomUser(t)
	testCases := []struct {
		name            string
		username        string
		password        string
		fullName        string
		email           string
		buildRepository func(repository *mockrepository.MockUserRepository)
		checkResponse   func(t *testing.T, res domain.User, err error)
	}{
		{
			name:     "OK",
			username: user.Username,
			password: user.HashedPassword,
			fullName: user.FullName,
			email:    user.Email,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, username, password, fullName, email string) (domain.User, error) {
					require.Equal(t, user.Username, username)
					require.Equal(t, user.FullName, fullName)
					require.Equal(t, user.Email, email)
					return user, nil
				})
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.NoError(t, err)
				require.Equal(t, user, res)
			},
		},
		{
			name:     "Failed",
			username: user.Username,
			password: user.HashedPassword,
			fullName: user.FullName,
			email:    user.Email,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, username, password, fullName, email string) (domain.User, error) {
					return domain.User{}, fmt.Errorf("duplicate key")
				})
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.Error(t, err)
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockUserRepository(ctrl)
			testCases[i].buildRepository(repository)
			userService := NewTestUserService(t, repository)
			res, err := userService.CreateUser(context.Background(), testCases[i].username, testCases[i].password, testCases[i].fullName, testCases[i].email)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	user := RandomUser(t)
	testCases := []struct {
		name            string
		username        string
		buildRepository func(repository *mockrepository.MockUserRepository)
		checkResponse   func(t *testing.T, res domain.User, err error)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(user, nil)
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.NoError(t, err)
				require.Equal(t, user, res)
			},
		},
		{
			name:     "NotFound",
			username: user.Username,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(domain.User{}, fmt.Errorf("not found"))
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.Error(t, err)
				require.Equal(t, "user not found", err.Error())
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockUserRepository(ctrl)
			testCases[i].buildRepository(repository)
			userService := NewTestUserService(t, repository)
			res, err := userService.GetUserByUsername(context.Background(), testCases[i].username)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestUpdateFullNameByUsername(t *testing.T) {
	user := RandomUser(t)
	newFullName := util.RandomFullName(6, 12)
	updatedUser := user
	updatedUser.FullName = newFullName

	testCases := []struct {
		name            string
		username        string
		fullName        string
		buildRepository func(repository *mockrepository.MockUserRepository)
		checkResponse   func(t *testing.T, res domain.User, err error)
	}{
		{
			name:     "OK",
			username: user.Username,
			fullName: newFullName,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().UpdateFullNameByUsername(gomock.Any(), gomock.Eq(user.Username), gomock.Eq(newFullName)).Return(updatedUser, nil)
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.NoError(t, err)
				require.Equal(t, updatedUser, res)
			},
		},
		{
			name:     "NotFound",
			username: user.Username,
			fullName: newFullName,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().UpdateFullNameByUsername(gomock.Any(), gomock.Eq(user.Username), gomock.Eq(newFullName)).Return(domain.User{}, fmt.Errorf("not found"))
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.Error(t, err)
				require.Equal(t, "user not found", err.Error())
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockUserRepository(ctrl)
			testCases[i].buildRepository(repository)
			userService := NewTestUserService(t, repository)
			res, err := userService.UpdateFullNameByUsername(context.Background(), testCases[i].username, testCases[i].fullName)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestUpdateUserNameByUsername(t *testing.T) {
	user := RandomUser(t)
	newUsername := util.RandomUserName(6, 12)
	updatedUser := user
	updatedUser.Username = newUsername

	testCases := []struct {
		name            string
		username        string
		newUsername     string
		buildRepository func(repository *mockrepository.MockUserRepository)
		checkResponse   func(t *testing.T, res domain.User, err error)
	}{
		{
			name:        "OK",
			username:    user.Username,
			newUsername: newUsername,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().UpdateUserNameByUsername(gomock.Any(), gomock.Eq(user.Username), gomock.Eq(newUsername)).Return(updatedUser, nil)
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.NoError(t, err)
				require.Equal(t, updatedUser, res)
			},
		},
		{
			name:        "NotFound",
			username:    user.Username,
			newUsername: newUsername,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().UpdateUserNameByUsername(gomock.Any(), gomock.Eq(user.Username), gomock.Eq(newUsername)).Return(domain.User{}, fmt.Errorf("not found"))
			},
			checkResponse: func(t *testing.T, res domain.User, err error) {
				require.Error(t, err)
				require.Equal(t, "user not found", err.Error())
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockUserRepository(ctrl)
			testCases[i].buildRepository(repository)
			userService := NewTestUserService(t, repository)
			res, err := userService.UpdateUserNameByUsername(context.Background(), testCases[i].username, testCases[i].newUsername)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestDeleteUserByUsername(t *testing.T) {
	user := RandomUser(t)
	testCases := []struct {
		name            string
		username        string
		buildRepository func(repository *mockrepository.MockUserRepository)
		checkResponse   func(t *testing.T, err error)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().DeleteUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "NotFound",
			username: user.Username,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().DeleteUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(fmt.Errorf("not found"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "user not found", err.Error())
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockUserRepository(ctrl)
			testCases[i].buildRepository(repository)
			userService := NewTestUserService(t, repository)
			err := userService.DeleteUserByUsername(context.Background(), testCases[i].username)
			testCases[i].checkResponse(t, err)
		})
	}
}
