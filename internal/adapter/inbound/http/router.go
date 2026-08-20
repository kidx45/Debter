package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(userAdapter *UserAdapter, accountAdapter *AccountAdapter, entryAdapter *EntryAdapter) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userAdapter.CreateUser)
			users.GET("/:username", userAdapter.GetUserByUsername)
			users.PUT("/:username/fullname", userAdapter.UpdateFullNameByUsername)
			users.PUT("/:username/username", userAdapter.UpdateUserNameByUsername)
			users.DELETE("/:username", userAdapter.DeleteUserByUsername)
		}

		accounts := v1.Group("/accounts")
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
