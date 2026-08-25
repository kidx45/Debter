package service

import (
	"context"
	"fmt"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type UserService struct {
	UserRepo outbound.UserRepository
}

func NewUserService(userRepo outbound.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, username, password, fullName, email string) (domain.User, error) {
	return u.UserRepo.CreateUser(ctx, username, password, fullName, email)
}

func (u *UserService) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := u.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *UserService) UpdateFullNameByUsername(ctx context.Context, username, fullName string) (domain.User, error) {
	user, err := u.UserRepo.UpdateFullNameByUsername(ctx, username, fullName)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *UserService) UpdateUserNameByUsername(ctx context.Context, username, newUsername string) (domain.User, error) {
	user, err := u.UserRepo.UpdateUserNameByUsername(ctx, username, newUsername)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *UserService) DeleteUserByUsername(ctx context.Context, username string) error {
	err := u.UserRepo.DeleteUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	return nil
}
