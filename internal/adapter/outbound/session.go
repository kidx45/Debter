package outbound

import (
	"context"
	"time"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type postgresSessionRepo struct {
	queries *db.Queries
}

func NewPostgresSessionRepository(queries *db.Queries) outbound.SessionRepository {
	return &postgresSessionRepo{queries: queries}
}

func (r *postgresSessionRepo) CreateSession(ctx context.Context, userID int64, refreshToken, userAgent, clientIp string, expiresAt time.Time) (domain.Session, error) {
	result, err := r.queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		return domain.Session{}, err
	}
	return dbSessionToDomain(result), nil
}

func (r *postgresSessionRepo) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	result, err := r.queries.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.Session{}, err
	}
	return dbSessionToDomain(result), nil
}

func (r *postgresSessionRepo) UpdateSessionRefreshToken(ctx context.Context, id int64, refreshToken string, expiresAt time.Time) (domain.Session, error) {
	result, err := r.queries.UpdateSessionRefreshToken(ctx, db.UpdateSessionRefreshTokenParams{
		ID:           id,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		return domain.Session{}, err
	}
	return dbSessionToDomain(result), nil
}

func dbSessionToDomain(s db.Session) domain.Session {
	return domain.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIp:     s.ClientIp,
		IsBlocked:    s.IsBlocked,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
	}
}
