package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/util"
	"github.com/kidx45/Debter/internal/util/token"
)

type LoginUserRequest struct {
	Username  string
	Password  string
	UserAgent string
	ClientIp  string
}

type LoginUserResult struct {
	User                  domain.User
	Session               domain.Session
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type RenewAccessTokenResult struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type AuthService struct {
	UserRepo             outbound.UserRepository
	SessionRepo          outbound.SessionRepository
	TokenMaker           token.TokenMaker
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewAuthService(userRepo outbound.UserRepository, sessionRepo outbound.SessionRepository, tokenMaker token.TokenMaker, accessTokenDuration, refreshTokenDuration time.Duration) *AuthService {
	return &AuthService{
		UserRepo:             userRepo,
		SessionRepo:          sessionRepo,
		TokenMaker:           tokenMaker,
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, req LoginUserRequest) (*LoginUserResult, error) {
	user, err := s.UserRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("incorrect username or password")
		}
		return nil, err
	}
	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return nil, fmt.Errorf("incorrect username or password")
	}

	accessToken, accessPayload, err := s.TokenMaker.CreateToken(user.ID, user.Username, s.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken := util.RandomString(32)
	expiresAt := time.Now().Add(s.RefreshTokenDuration)
	session, err := s.SessionRepo.CreateSession(ctx, user.ID, refreshToken, req.UserAgent, req.ClientIp, expiresAt)
	if err != nil {
		return nil, err
	}

	return &LoginUserResult{
		User:                  user,
		Session:               session,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *AuthService) RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenResult, error) {
	session, err := s.SessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, err
	}
	if session.IsBlocked {
		return nil, fmt.Errorf("blocked session")
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("expired session")
	}

	user, err := s.UserRepo.GetUser(ctx, session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	newRefreshToken := util.RandomString(32)
	newExpiresAt := time.Now().Add(s.RefreshTokenDuration)
	updatedSession, err := s.SessionRepo.UpdateSessionRefreshToken(ctx, session.ID, newRefreshToken, newExpiresAt)
	if err != nil {
		return nil, err
	}

	accessToken, accessPayload, err := s.TokenMaker.CreateToken(user.ID, user.Username, s.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	return &RenewAccessTokenResult{
		AccessToken:           accessToken,
		RefreshToken:          updatedSession.RefreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: updatedSession.ExpiresAt,
	}, nil
}
