package outbound

import (
	"context"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/util"
)

type postgresUserRepo struct {
	queries *db.Queries
}

func NewPostgresUserRepository(queries *db.Queries) outbound.UserRepository {
	return &postgresUserRepo{queries: queries}
}

func (r *postgresUserRepo) CreateUser(ctx context.Context, username, hashedPassword, fullName, email string) (domain.User, error) {
	hashed, err := util.HashPassword(hashedPassword)
	if err != nil {
		return domain.User{}, err
	}
	result, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:       username,
		HashedPassword: hashed,
		FullName:       fullName,
		Email:          email,
	})
	if err != nil {
		return domain.User{}, err
	}
	return dbUserToDomain(result), nil
}

func (r *postgresUserRepo) DeleteUserByUsername(ctx context.Context, username string) error {
	return r.queries.DeleteUserByUsername(ctx, username)
}

func (r *postgresUserRepo) GetUser(ctx context.Context, userID int64) (domain.User, error) {
	result, err := r.queries.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return dbUserToDomain(result), nil
}

func (r *postgresUserRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	result, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return dbUserToDomain(result), nil
}

func (r *postgresUserRepo) UpdateFullNameByUsername(ctx context.Context, username, fullName string) (domain.User, error) {
	result, err := r.queries.UpdateFullNameByUsername(ctx, db.UpdateFullNameByUsernameParams{
		FullName: fullName,
		Username: username,
	})
	if err != nil {
		return domain.User{}, err
	}
	return dbUserToDomain(result), nil
}

func (r *postgresUserRepo) UpdateUserNameByUsername(ctx context.Context, username, newUsername string) (domain.User, error) {
	result, err := r.queries.UpdateUserNameByUsername(ctx, db.UpdateUserNameByUsernameParams{
		Username:    username,
		NewUsername: newUsername,
	})
	if err != nil {
		return domain.User{}, err
	}
	return dbUserToDomain(result), nil
}

func dbUserToDomain(u db.User) domain.User {
	return domain.User{
		ID:             u.ID,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		FullName:       u.FullName,
		Email:          u.Email,
		CreatedAt:      u.CreatedAt,
	}
}
