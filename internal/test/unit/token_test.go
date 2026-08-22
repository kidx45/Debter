package unit

import (
	"testing"
	"time"

	"github.com/kidx45/Debter/internal/util"
	"github.com/kidx45/Debter/internal/util/token"
	"github.com/stretchr/testify/require"
)

func randomSymmetricKey() string {
	return util.RandomString(32)
}

func TestTokenMaker(t *testing.T) {
	key := randomSymmetricKey()

	makers := []struct {
		name  string
		maker func(t *testing.T) token.TokenMaker
	}{
		{
			name: "Paseto",
			maker: func(t *testing.T) token.TokenMaker {
				maker, err := token.NewPasetoMaker(key)
				require.NoError(t, err)
				return maker
			},
		},
		{
			name: "JWT",
			maker: func(t *testing.T) token.TokenMaker {
				maker, err := token.NewJWTMaker(key)
				require.NoError(t, err)
				return maker
			},
		},
	}

	for _, m := range makers {
		t.Run(m.name, func(t *testing.T) {
			maker := m.maker(t)
			username := util.RandomUserName(6, 12)
			userID := util.RandomNumber(1, 1000)

			t.Run("OK", func(t *testing.T) {
				duration := time.Minute
				issuedAt := time.Now()
				expiredAt := issuedAt.Add(duration)

				res, payload, err := maker.CreateToken(userID, username, duration)
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.NotEmpty(t, payload)

				payload, err = maker.VerifyToken(res)
				require.NoError(t, err)
				require.NotEmpty(t, payload)

				require.Equal(t, userID, payload.UserID)
				require.Equal(t, username, payload.Username)
				require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
				require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
			})

			t.Run("ExpiredToken", func(t *testing.T) {
				res, _, err := maker.CreateToken(userID, username, -time.Minute)
				require.NoError(t, err)
				require.NotEmpty(t, res)

				payload, err := maker.VerifyToken(res)
				require.Error(t, err)
				require.Equal(t, token.ErrExpiredToken, err)
				require.Nil(t, payload)
			})

			t.Run("InvalidToken", func(t *testing.T) {
				payload, err := maker.VerifyToken(util.RandomString(32))
				require.Error(t, err)
				require.ErrorIs(t, err, token.ErrInvalidToken)
				require.Nil(t, payload)
			})
		})
	}
}

func TestNewPasetoMakerInvalidKeySize(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(util.RandomNumber(1, 31)))
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestNewJWTMakerInvalidKeySize(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(util.RandomNumber(1, 31)))
	require.Error(t, err)
	require.Nil(t, maker)
}
