package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Debter/internal/service"
)

type AuthAdapter struct {
	AuthService *service.AuthService
}

func NewAuthAdapter(AuthService service.AuthService) *AuthAdapter {
	return &AuthAdapter{
		AuthService: &AuthService,
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	User                  userResponse `json:"user"`
	SessionID             int64        `json:"sessionId"`
	AccessToken           string       `json:"accessToken"`
	RefreshToken          string       `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time    `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time    `json:"refreshTokenExpiresAt"`
}

func (a *AuthAdapter) Login(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.AuthService.LoginUser(ctx.Request.Context(), service.LoginUserRequest{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp:  ctx.ClientIP(),
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		User:                  newUserResponse(res.User),
		SessionID:             res.Session.ID,
		AccessToken:           res.AccessToken,
		RefreshToken:          res.RefreshToken,
		AccessTokenExpiresAt:  res.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: res.RefreshTokenExpiresAt,
	})
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

func (a *AuthAdapter) RenewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.AuthService.RenewAccessToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:           res.AccessToken,
		RefreshToken:          res.RefreshToken,
		AccessTokenExpiresAt:  res.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: res.RefreshTokenExpiresAt,
	})
}
