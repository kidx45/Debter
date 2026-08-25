package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Debter/internal/adapter/inbound/middleware"
	"github.com/kidx45/Debter/internal/service"
)

type AccountAdapter struct {
	AccountService *service.AccountService
}

func NewAccountAdapter(accountService service.AccountService) *AccountAdapter {
	return &AccountAdapter{
		AccountService: &accountService,
	}
}

func (a *AccountAdapter) GetAccountsByUserId(ctx *gin.Context) {
	payload := middleware.GetAuthPayload(ctx)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload is not provided"})
		return
	}

	userID, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if payload.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "account does not belong to the authenticated user"})
		return
	}

	res, err := a.AccountService.GetAccountsByUserId(ctx.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
