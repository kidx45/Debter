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
	"github.com/kidx45/Debter/internal/util/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func RandomUserWithPassword(t *testing.T) (db.User, string) {
	password := util.RandomPassword(8, 16)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:             util.RandomNumber(1, 1000),
		Username:       util.RandomUserName(6, 12),
		HashedPassword: hashedPassword,
		FullName:       util.RandomFullName(6, 12),
		Email:          util.RandomEmail(6, 12),
		CreatedAt:      time.Now(),
	}
	return user, password
}

func NewTestAuthService(t *testing.T, userRepo outbound.UserRepository, sessionRepo outbound.SessionRepository, maker token.TokenMaker) *service.AuthService {
	return service.NewAuthService(userRepo, sessionRepo, maker, time.Minute, 24*time.Hour)
}

func NewTestTokenMaker(t *testing.T) token.TokenMaker {
	maker, err := token.NewPasetoMaker(randomSymmetricKey())
	require.NoError(t, err)
	return maker
}

func TestLoginUser(t *testing.T) {
	user, password := RandomUserWithPassword(t)
	maker := NewTestTokenMaker(t)

	testCases := []struct {
		name          string
		req           service.LoginUserRequest
		buildStubs    func(userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository)
		checkResponse func(t *testing.T, res *service.LoginUserResult, err error)
	}{
		{
			name: "OK",
			req: service.LoginUserRequest{
				Username:  user.Username,
				Password:  password,
				UserAgent: util.RandomString(10),
				ClientIp:  "127.0.0.1",
			},
			buildStubs: func(userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				userRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(user, nil)
				sessionRepo.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
					require.Equal(t, user.ID, arg.UserID)
					require.NotEmpty(t, arg.RefreshToken)
					require.True(t, arg.ExpiresAt.After(time.Now()))
					return db.Session{
						ID:           util.RandomNumber(1, 1000),
						UserID:       arg.UserID,
						RefreshToken: arg.RefreshToken,
						UserAgent:    arg.UserAgent,
						ClientIp:     arg.ClientIp,
						IsBlocked:    false,
						ExpiresAt:    arg.ExpiresAt,
						CreatedAt:    time.Now(),
					}, nil
				})
			},
			checkResponse: func(t *testing.T, res *service.LoginUserResult, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)

				payload, verifyErr := maker.VerifyToken(res.AccessToken)
				require.NoError(t, verifyErr)
				require.Equal(t, user.ID, payload.UserID)
				require.Equal(t, user.Username, payload.Username)
				require.WithinDuration(t, time.Now().Add(time.Minute), res.AccessTokenExpiresAt, 2*time.Second)
				require.WithinDuration(t, time.Now().Add(24*time.Hour), res.RefreshTokenExpiresAt, 2*time.Second)
				require.Equal(t, user.ID, res.User.ID)
				require.Equal(t, user.ID, res.Session.UserID)
				require.Equal(t, res.RefreshToken, res.Session.RefreshToken)
			},
		},
		{
			name: "WrongPassword",
			req: service.LoginUserRequest{
				Username: user.Username,
				Password: util.RandomPassword(8, 16),
			},
			buildStubs: func(userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				userRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Return(user, nil)
			},
			checkResponse: func(t *testing.T, res *service.LoginUserResult, err error) {
				require.Error(t, err)
				require.Equal(t, "incorrect username or password", err.Error())
				require.Nil(t, res)
			},
		},
		{
			name: "UnknownUser",
			req: service.LoginUserRequest{
				Username: util.RandomUserName(6, 12),
				Password: password,
			},
			buildStubs: func(userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				userRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *service.LoginUserResult, err error) {
				require.Error(t, err)
				require.Equal(t, "incorrect username or password", err.Error())
				require.Nil(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mockrepository.NewMockUserRepository(ctrl)
			sessionRepo := mockrepository.NewMockSessionRepository(ctrl)
			testCases[i].buildStubs(userRepo, sessionRepo)
			authService := NewTestAuthService(t, userRepo, sessionRepo, maker)
			res, err := authService.LoginUser(context.Background(), testCases[i].req)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestRenewAccessToken(t *testing.T) {
	user, _ := RandomUserWithPassword(t)
	maker := NewTestTokenMaker(t)

	randomSession := func() db.Session {
		return db.Session{
			ID:           util.RandomNumber(1, 1000),
			UserID:       user.ID,
			RefreshToken: util.RandomString(32),
			UserAgent:    util.RandomString(10),
			ClientIp:     fmt.Sprintf("192.168.%d.%d", util.RandomNumber(0, 255), util.RandomNumber(0, 255)),
			IsBlocked:    false,
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
		}
	}

	testCases := []struct {
		name          string
		buildStubs    func(session *db.Session, userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository)
		checkResponse func(t *testing.T, session *db.Session, res *service.RenewAccessTokenResult, err error)
	}{
		{
			name: "OK",
			buildStubs: func(session *db.Session, userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				sessionRepo.EXPECT().GetSessionByRefreshToken(gomock.Any(), gomock.Eq(session.RefreshToken)).Return(*session, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Return(user, nil)
				sessionRepo.EXPECT().UpdateSessionRefreshToken(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, arg db.UpdateSessionRefreshTokenParams) (db.Session, error) {
					require.Equal(t, session.ID, arg.ID)
					require.NotEqual(t, session.RefreshToken, arg.RefreshToken)
					require.True(t, arg.ExpiresAt.After(time.Now()))
					return db.Session{
						ID:           session.ID,
						UserID:       session.UserID,
						RefreshToken: arg.RefreshToken,
						UserAgent:    session.UserAgent,
						ClientIp:     session.ClientIp,
						IsBlocked:    false,
						ExpiresAt:    arg.ExpiresAt,
						CreatedAt:    session.CreatedAt,
					}, nil
				})
			},
			checkResponse: func(t *testing.T, session *db.Session, res *service.RenewAccessTokenResult, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.AccessToken)
				require.NotEqual(t, session.RefreshToken, res.RefreshToken)
				require.WithinDuration(t, time.Now().Add(time.Minute), res.AccessTokenExpiresAt, 2*time.Second)
				require.WithinDuration(t, time.Now().Add(24*time.Hour), res.RefreshTokenExpiresAt, 2*time.Second)
			},
		},
		{
			name: "InvalidSession",
			buildStubs: func(session *db.Session, userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				sessionRepo.EXPECT().GetSessionByRefreshToken(gomock.Any(), gomock.Eq(session.RefreshToken)).Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, session *db.Session, res *service.RenewAccessTokenResult, err error) {
				require.Error(t, err)
				require.Equal(t, "invalid refresh token", err.Error())
				require.Nil(t, res)
			},
		},
		{
			name: "BlockedSession",
			buildStubs: func(session *db.Session, userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				session.IsBlocked = true
				sessionRepo.EXPECT().GetSessionByRefreshToken(gomock.Any(), gomock.Eq(session.RefreshToken)).Return(*session, nil)
			},
			checkResponse: func(t *testing.T, session *db.Session, res *service.RenewAccessTokenResult, err error) {
				require.Error(t, err)
				require.Equal(t, "blocked session", err.Error())
				require.Nil(t, res)
			},
		},
		{
			name: "ExpiredSession",
			buildStubs: func(session *db.Session, userRepo *mockrepository.MockUserRepository, sessionRepo *mockrepository.MockSessionRepository) {
				session.ExpiresAt = time.Now().Add(-time.Hour)
				sessionRepo.EXPECT().GetSessionByRefreshToken(gomock.Any(), gomock.Eq(session.RefreshToken)).Return(*session, nil)
			},
			checkResponse: func(t *testing.T, session *db.Session, res *service.RenewAccessTokenResult, err error) {
				require.Error(t, err)
				require.Equal(t, "expired session", err.Error())
				require.Nil(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			session := randomSession()
			userRepo := mockrepository.NewMockUserRepository(ctrl)
			sessionRepo := mockrepository.NewMockSessionRepository(ctrl)
			testCases[i].buildStubs(&session, userRepo, sessionRepo)
			authService := NewTestAuthService(t, userRepo, sessionRepo, maker)
			res, err := authService.RenewAccessToken(context.Background(), session.RefreshToken)
			testCases[i].checkResponse(t, &session, res, err)
		})
	}
}
