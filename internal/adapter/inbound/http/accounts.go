package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Debter/internal/service"
)

type AccountAdapter struct {
	AccountService *service.AccountService
}

func NewAccountAdapter(AccountService service.AccountService) *AccountAdapter {
	return &AccountAdapter{
		AccountService: &AccountService,
	}
}

func (a *AccountAdapter) GetAccountsByUserId(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
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
