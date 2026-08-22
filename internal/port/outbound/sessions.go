package outbound

import (
	"context"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (db.Session, error)
	UpdateSessionRefreshToken(ctx context.Context, arg db.UpdateSessionRefreshTokenParams) (db.Session, error)
}
