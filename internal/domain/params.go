package domain

import "time"

type CreateUserParams struct {
	Username       string
	HashedPassword string
	FullName       string
	Email          string
}

type CreateSessionParams struct {
	UserID       int64
	RefreshToken string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
}

type GetEntriesByAccountIdParams struct {
	AccountID int64
	UserID    int64
}

type FilterEntriesByDateParams struct {
	AccountID   int64
	UserID      int64
	CreatedAt   time.Time
	CreatedAt_2 time.Time
}

type GetEntriesByCategoryAndTypeParams struct {
	AccountID int64
	UserID    int64
	Type      string
}

type UpdateFullNameByUsernameParams struct {
	Username string
	FullName string
}

type UpdateUserNameByUsernameParams struct {
	Username    string
	NewUsername string
}

type UpdateSessionRefreshTokenParams struct {
	ID           int64
	RefreshToken string
	ExpiresAt    time.Time
}
