package http

import "github.com/kidx45/Debter/internal/service"

type UserAdapter struct {
	UserService *service.UserService
}

func NewUserAdapter (UserService service.UserService) *UserAdapter {
	return &UserAdapter{
		UserService: &UserService,
	}
}
