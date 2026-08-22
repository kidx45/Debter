package unit

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/service"
	mockrepository "github.com/kidx45/Debter/internal/test/mock/repository"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func RandomUser(t *testing.T) db.User {
	return db.User{
		ID:             util.RandomNumber(1, 1000),
		FullName:       util.RandomFullName(6, 12),
		Username:       util.RandomUserName(6, 12),
		Email:          util.RandomEmail(6, 12),
		HashedPassword: util.RandomPassword(8, 16),
		CreatedAt:      time.Now(),
	}
}

func NewTestUserService(t *testing.T, DB outbound.UserRepository) *service.UserService {
	return service.NewUserService(DB)
}

func TestCreateUser(t *testing.T) {
	user := RandomUser(t)
	testCases := []struct {
		name            string
		req             db.User
		buildRepository func(repository *mockrepository.MockUserRepository)
		checkResponse   func(t *testing.T, res db.User, err error)
	}{
		{
			name: "OK",
			req:  user,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
					require.NoError(t, util.CheckPassword(user.HashedPassword, arg.HashedPassword))
					require.NotEqual(t, user.HashedPassword, arg.HashedPassword)
					require.Equal(t, user.Username, arg.Username)
					require.Equal(t, user.FullName, arg.FullName)
					require.Equal(t, user.Email, arg.Email)
					return user, nil
				})
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
				require.NoError(t, err)
				require.Equal(t, user, res)
			},
		},
		{
			name: "Failed",
			req:  user,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
					require.NoError(t, util.CheckPassword(user.HashedPassword, arg.HashedPassword))
					return db.User{}, fmt.Errorf("duplicate key")
				})
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
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
			UserService := NewTestUserService(t, repository)
			res, err := UserService.CreateUser(context.Background(), testCases[i].req)
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
		checkResponse   func(t *testing.T, res db.User, err error)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(user, nil)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
				require.NoError(t, err)
				require.Equal(t, user, res)
			},
		},
		{
			name:     "NotFound",
			username: user.Username,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				repository.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
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
			UserService := NewTestUserService(t, repository)
			res, err := UserService.GetUserByUsername(context.Background(), testCases[i].username)
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
		checkResponse   func(t *testing.T, res db.User, err error)
	}{
		{
			name:     "OK",
			username: user.Username,
			fullName: newFullName,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				arg := db.UpdateFullNameByUsernameParams{
					Username: user.Username,
					FullName: newFullName,
				}
				repository.EXPECT().UpdateFullNameByUsername(gomock.Any(), gomock.Eq(arg)).Return(updatedUser, nil)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
				require.NoError(t, err)
				require.Equal(t, updatedUser, res)
			},
		},
		{
			name:     "NotFound",
			username: user.Username,
			fullName: newFullName,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				arg := db.UpdateFullNameByUsernameParams{
					Username: user.Username,
					FullName: newFullName,
				}
				repository.EXPECT().UpdateFullNameByUsername(gomock.Any(), gomock.Eq(arg)).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
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
			UserService := NewTestUserService(t, repository)
			res, err := UserService.UpdateFullNameByUsername(context.Background(), testCases[i].username, testCases[i].fullName)
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
		checkResponse   func(t *testing.T, res db.User, err error)
	}{
		{
			name:        "OK",
			username:    user.Username,
			newUsername: newUsername,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				arg := db.UpdateUserNameByUsernameParams{
					Username:    user.Username,
					NewUsername: newUsername,
				}
				repository.EXPECT().UpdateUserNameByUsername(gomock.Any(), gomock.Eq(arg)).Return(updatedUser, nil)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
				require.NoError(t, err)
				require.Equal(t, updatedUser, res)
			},
		},
		{
			name:        "NotFound",
			username:    user.Username,
			newUsername: newUsername,
			buildRepository: func(repository *mockrepository.MockUserRepository) {
				arg := db.UpdateUserNameByUsernameParams{
					Username:    user.Username,
					NewUsername: newUsername,
				}
				repository.EXPECT().UpdateUserNameByUsername(gomock.Any(), gomock.Eq(arg)).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
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
			UserService := NewTestUserService(t, repository)
			res, err := UserService.UpdateUserNameByUsername(context.Background(), testCases[i].username, testCases[i].newUsername)
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
				repository.EXPECT().DeleteUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(sql.ErrNoRows)
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
			UserService := NewTestUserService(t, repository)
			err := UserService.DeleteUserByUsername(context.Background(), testCases[i].username)
			testCases[i].checkResponse(t, err)
		})
	}
}
