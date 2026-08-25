package outbound

import (
	"context"
	"time"

	"github.com/kidx45/Debter/internal/domain"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, userID int64, refreshToken, userAgent, clientIp string, expiresAt time.Time) (domain.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error)
	UpdateSessionRefreshToken(ctx context.Context, id int64, refreshToken string, expiresAt time.Time) (domain.Session, error)
}
