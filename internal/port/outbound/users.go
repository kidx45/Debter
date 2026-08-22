package outbound

import (
	"context"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUserByUsername(ctx context.Context, username string) error
	GetUser(ctx context.Context, userID int64) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	UpdateFullNameByUsername(ctx context.Context, arg db.UpdateFullNameByUsernameParams) (db.User, error)
	UpdateUserNameByUsername(ctx context.Context, arg db.UpdateUserNameByUsernameParams) (db.User, error)
}

