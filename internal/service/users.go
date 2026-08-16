package service

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
)
type UserService struct {
	DB outbound.UserRepository
}

func NewUserService(DB outbound.UserRepository) *UserService {
	return &UserService{
		DB: DB,
	}
}

func (u *UserService) CreateUser(ctx context.Context, req db.CreateUserParams) (db.User, error) {
	res, err := u.DB.CreateUser(ctx,req)
	if err != nil {
		return db.User{}, err
	}
}