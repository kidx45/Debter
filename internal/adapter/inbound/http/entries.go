package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Debter/internal/adapter/inbound/middleware"
	"github.com/kidx45/Debter/internal/service"
)

type EntryAdapter struct {
	EntryService *service.EntryService
}

func NewEntryAdapter(entryService service.EntryService) *EntryAdapter {
	return &EntryAdapter{
		EntryService: &entryService,
	}
}

func (a *EntryAdapter) GetEntriesByAccountId(ctx *gin.Context) {
	payload := middleware.GetAuthPayload(ctx)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload is not provided"})
		return
	}

	accountID, err := strconv.ParseInt(ctx.Param("accountId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	res, err := a.EntryService.GetEntriesByAccountId(ctx.Request.Context(), accountID, payload.UserID)
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

	accountID, err := strconv.ParseInt(ctx.Param("accountId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	from, err := time.Parse(time.RFC3339, ctx.Query("from"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date, use RFC3339 format"})
		return
	}

	to, err := time.Parse(time.RFC3339, ctx.Query("to"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date, use RFC3339 format"})
		return
	}

	res, err := a.EntryService.FilterEntriesByDate(ctx.Request.Context(), accountID, payload.UserID, from, to)
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

	accountID, err := strconv.ParseInt(ctx.Param("accountId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	entryType := ctx.Query("type")
	if entryType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'type' query parameter is required"})
		return
	}

	res, err := a.EntryService.GetEntriesByCategoryAndType(ctx.Request.Context(), accountID, payload.UserID, entryType)
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
