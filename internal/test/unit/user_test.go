package unit

import (
	"context"
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

func RandomUser (t* testing.T) db.User {
	return db.User{
		ID: util.RandomNumber(1, 1000),
		FullName: util.RandomFullName(6, 12),
		Username: util.RandomUserName(6, 12),
		Email: util.RandomEmail(6, 12),
		HashedPassword: util.RandomPassword(8, 16),
		CreatedAt: time.Now(),
	}
}

func NewTestUserService(t *testing.T, DB outbound.UserRepository) *service.UserService {
	return service.NewUserService(DB)
}

func TestCreateUser(t *testing.T) {
	user := RandomUser(t)
	testCases := []struct {
		name string
		req db.User
		buildRepository func(repository *mockrepository.MockRepository)
		checkResponse func(t *testing.T, res db.User, err error)
	} {
		{
			name: "OK",
			req: user,
			buildRepository: func(repository *mockrepository.MockRepository) {
				arg := db.CreateUserParams{
					Username: user.Username,
					HashedPassword: user.HashedPassword,
					FullName: user.FullName,
					Email: user.Email,
				}
				repository.EXPECT().CreateUser(gomock.Any(), gomock.Eq(arg)).Return(user, nil)
			},
			checkResponse: func(t *testing.T, res db.User, err error) {
				require.NoError(t, err)
				require.Equal(t, user, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name,func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			repository := mockrepository.NewMockRepository(ctrl)
			testCases[i].buildRepository(repository)
			UserService := NewTestUserService(t, repository)
			res, err := UserService.CreateUser(context.Background(), testCases[i].req)
			testCases[i].checkResponse(t, res, err)
		})
	}
}