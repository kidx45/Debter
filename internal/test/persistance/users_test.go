package persistance

import (
	"context"
	"testing"

	"github.com/kidx45/Debter/internal/adapter/outbound/db"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) db.User {
	arg := db.CreateUserParams{
		Username: util.RandomString(6),
		FullName: util.RandomString(6),
		Email: util.RandomString(6),
		HashedPassword: util.RandomString(6),
	}

	res, err := testdb.CreateUser(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,res.ID)
	require.NotEmpty(t,res.Username)
	require.NotEmpty(t,res.FullName)
	require.NotEmpty(t,res.Email)
	require.NotEmpty(t,res.CreatedAt)
	return res
}
