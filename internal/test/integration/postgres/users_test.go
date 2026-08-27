package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) (db.User, *db.Queries) {
	query := TxQueriesTest(t)
	arg := db.CreateUserParams{
		Username:       util.RandomString(6),
		FullName:       util.RandomString(6),
		Email:          util.RandomString(6),
		HashedPassword: util.RandomString(6),
	}

	res, err := query.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, arg.Username, res.Username)
	require.Equal(t, arg.HashedPassword, res.HashedPassword)
	require.Equal(t, arg.Email, res.Email)
	require.Equal(t, arg.FullName, res.FullName)
	return res, query
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUserByUsername(t *testing.T) {
	user, query := CreateRandomUser(t)

	res, err := query.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, user, res)
}

func TestUpdateFullNameByUsername(t *testing.T) {
	user, query := CreateRandomUser(t)
	arg := db.UpdateFullNameByUsernameParams{
		FullName: util.RandomString(6),
		Username: user.Username,
	}
	res, err := query.UpdateFullNameByUsername(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, user.FullName, res.FullName)
	require.Equal(t, arg.FullName, res.FullName)
	require.Equal(t, user.ID, res.ID)
	require.Equal(t, user.Email, res.Email)
	require.Equal(t, user.HashedPassword, res.HashedPassword)
	require.Equal(t, user.CreatedAt, res.CreatedAt)
	require.Equal(t, user.Username, res.Username)
}

func TestUpdateUserNameByUsername(t *testing.T) {
	user, query := CreateRandomUser(t)
	arg := db.UpdateUserNameByUsernameParams{
		Username:    user.Username,
		NewUsername: util.RandomString(6),
	}
	res, err := query.UpdateUserNameByUsername(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, user.Username, res.Username)
	require.Equal(t, arg.NewUsername, res.Username)
	require.Equal(t, user.ID, res.ID)
	require.Equal(t, user.Email, res.Email)
	require.Equal(t, user.HashedPassword, res.HashedPassword)
	require.Equal(t, user.CreatedAt, res.CreatedAt)
	require.Equal(t, user.FullName, res.FullName)
}

func TestDeleteUserByUsername(t *testing.T) {
	user, query := CreateRandomUser(t)
	err := query.DeleteUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	res, err := query.GetUserByUsername(context.Background(), user.Username)
	require.ErrorIs(t, sql.ErrNoRows, err)
	require.Empty(t, res)
}
