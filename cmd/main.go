package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	httpAdapter "github.com/kidx45/Debter/internal/adapter/inbound/http"
	"github.com/kidx45/Debter/internal/adapter/outbound"
	"github.com/kidx45/Debter/internal/service"
	"github.com/kidx45/Debter/internal/util"
)

func main() {
	config, err := util.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading the Env Configuration: %s", err)
	}

	conn, err := sql.Open(config.DB_DRIVER_NAME, config.DB_URL)
	if err != nil {
		log.Fatalf("Error setting up db connection: %s", err)
	}
	defer conn.Close()

	userRepo := outbound.NewPostgresUserRepository(conn)
	accountRepo := outbound.NewPostgresAccountRepository(conn)
	entryRepo := outbound.NewPostgresEntryRepository(conn)

	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	entryService := service.NewEntryService(entryRepo)

	userAdapter := httpAdapter.NewUserAdapter(*userService)
	accountAdapter := httpAdapter.NewAccountAdapter(*accountService)
	entryAdapter := httpAdapter.NewEntryAdapter(*entryService)

	router := httpAdapter.NewRouter(userAdapter, accountAdapter, entryAdapter)

	fmt.Println("Server starting on port " + config.PORT)
	if err := router.Run(":" + config.PORT); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
