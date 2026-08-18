package service

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/lib/pq"
)
type UserService struct {
	DB outbound.UserRepository
}

func NewUserService(DB outbound.UserRepository) *UserService {
	return &UserService{
		DB: DB,
	}
}

func (u *UserService) CreateUser(ctx context.Context, req db.User) (db.User, error) {
	// TODO: Fetch users username quickly for the frontend
	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: req.HashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}
	res, err := u.DB.CreateUser(ctx, arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name(){
				case "unique_violation":
					return db.User{}, fmt.Errorf("Email or Username has been already used")
			}
		}
		return db.User{}, err
	}
	return res,err
}

func (u *UserService) GetUserByUsername(ctx context.Context, Username string) (db.User, error) {
	res, err := u.DB.GetUserByUsername(ctx, Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found")
		}
		return db.User{}, err
	}
	return res, nil
}

func (u *UserService) UpdateFullNameByUsername(ctx context.Context, Username, FullName string) (db.User, error) {
	arg := db.UpdateFullNameByUsernameParams{
		Username: Username,
		FullName: FullName,
	}
	res, err := u.DB.UpdateFullNameByUsername(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found")
		}
		return db.User{}, err
	}
	return res, nil
}

func (u *UserService) UpdateUserNameByUsername(ctx context.Context, Username, NewUsername string) (db.User, error) {
	arg := db.UpdateUserNameByUsernameParams{
		Username: Username,
		NewUsername: NewUsername,
	}
	res, err := u.DB.UpdateUserNameByUsername(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found")
		}
		return db.User{}, err
	}
	return res, nil
}

func (u *UserService) DeleteUserByUsername(ctx context.Context, Username string) error {
	err := u.DB.DeleteUserByUsername(ctx, Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return err
	}
	return nil
}