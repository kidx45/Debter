package outbound

import (
	"context"

	"github.com/kidx45/Debter/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, hashedPassword, fullName, email string) (domain.User, error)
	DeleteUserByUsername(ctx context.Context, username string) error
	GetUser(ctx context.Context, userID int64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	UpdateFullNameByUsername(ctx context.Context, username, fullName string) (domain.User, error)
	UpdateUserNameByUsername(ctx context.Context, username, newUsername string) (domain.User, error)
}
