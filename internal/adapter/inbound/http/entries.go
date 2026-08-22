package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Debter/internal/adapter/inbound/middleware"
	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/service"
)

type EntryAdapter struct {
	EntryService *service.EntryService
}

func NewEntryAdapter(EntryService service.EntryService) *EntryAdapter {
	return &EntryAdapter{
		EntryService: &EntryService,
	}
}

func (a *EntryAdapter) GetEntriesByAccountId(ctx *gin.Context) {
	payload := middleware.GetAuthPayload(ctx)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload is not provided"})
		return
	}

	accountIDStr := ctx.Param("accountId")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	arg := db.GetEntriesByAccountIdParams{
		AccountID: accountID,
		UserID:    payload.UserID,
	}

	res, err := a.EntryService.GetEntriesByAccountId(ctx.Request.Context(), arg)
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

func (a *EntryAdapter) FilterEntriesByDate(ctx *gin.Context) {
	payload := middleware.GetAuthPayload(ctx)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload is not provided"})
		return
	}

	accountIDStr := ctx.Param("accountId")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date, use RFC3339 format"})
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date, use RFC3339 format"})
		return
	}

	arg := db.FilterEntriesByDateParams{
		AccountID:   accountID,
		UserID:      payload.UserID,
		CreatedAt:   from,
		CreatedAt_2: to,
	}

	res, err := a.EntryService.FilterEntriesByDate(ctx.Request.Context(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (a *EntryAdapter) GetEntriesByCategoryAndType(ctx *gin.Context) {
	payload := middleware.GetAuthPayload(ctx)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload is not provided"})
		return
	}

	accountIDStr := ctx.Param("accountId")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	entryType := ctx.Query("type")
	if entryType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'type' query parameter is required"})
		return
	}

	arg := db.GetEntriesByCategoryAndTypeParams{
		AccountID: accountID,
		UserID:    payload.UserID,
		Type:      entryType,
	}

	res, err := a.EntryService.GetEntriesByCategoryAndType(ctx.Request.Context(), arg)
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
