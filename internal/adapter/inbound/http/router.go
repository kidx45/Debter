package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kidx45/Debter/internal/adapter/inbound/middleware"
	"github.com/kidx45/Debter/internal/util/token"
)

func NewRouter(tokenMaker token.TokenMaker, userAdapter *UserAdapter, accountAdapter *AccountAdapter, entryAdapter *EntryAdapter, authAdapter *AuthAdapter, healthAdapter *HealthAdapter) *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", healthAdapter.Healthz)
	router.GET("/readyz", healthAdapter.Readyz)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userAdapter.CreateUser)
			users.POST("/login", authAdapter.Login)
			users.POST("/renew_access", authAdapter.RenewAccessToken)

			protected := users.Group("", middleware.AuthMiddleware(tokenMaker))
			{
				protected.GET("/:username", userAdapter.GetUserByUsername)
				protected.PUT("/:username/fullname", userAdapter.UpdateFullNameByUsername)
				protected.PUT("/:username/username", userAdapter.UpdateUserNameByUsername)
				protected.DELETE("/:username", userAdapter.DeleteUserByUsername)
			}
		}

		accounts := v1.Group("/accounts", middleware.AuthMiddleware(tokenMaker))
		{
			accounts.GET("/user/:userId", accountAdapter.GetAccountsByUserId)

			entries := accounts.Group("/:accountId/entries")
			{
				entries.GET("", entryAdapter.GetEntriesByAccountId)
				entries.GET("/filter", entryAdapter.FilterEntriesByDate)
				entries.GET("/summary", entryAdapter.GetEntriesByCategoryAndType)
			}
		}
	}

	return router
}
