package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	httpAdapter "github.com/kidx45/Debter/internal/adapter/inbound/http"
	"github.com/stretchr/testify/require"
)

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/healthz", httpAdapter.NewHealthAdapter(func(ctx context.Context) error { return nil }).Healthz)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestReadyz(t *testing.T) {
	testCases := []struct {
		name         string
		ping         func(ctx context.Context) error
		expectedCode int
	}{
		{
			name:         "database is ready",
			ping:         func(ctx context.Context) error { return nil },
			expectedCode: http.StatusOK,
		},
		{
			name:         "database is not ready",
			ping:         func(ctx context.Context) error { return context.DeadlineExceeded },
			expectedCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.GET("/readyz", httpAdapter.NewHealthAdapter(tc.ping).Readyz)

			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
