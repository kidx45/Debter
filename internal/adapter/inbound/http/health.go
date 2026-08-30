package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthAdapter struct {
	ping func(ctx context.Context) error
}

func NewHealthAdapter(ping func(ctx context.Context) error) *HealthAdapter {
	return &HealthAdapter{ping: ping}
}

func (h *HealthAdapter) Healthz(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (h *HealthAdapter) Readyz(ctx *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.ping(timeoutCtx); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "database is not ready"})
		return
	}

	ctx.Status(http.StatusOK)
}
