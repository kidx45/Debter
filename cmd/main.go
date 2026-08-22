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
	"github.com/kidx45/Debter/internal/util/token"
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

	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		log.Fatalf("Error setting up token maker: %s", err)
	}

	userRepo := outbound.NewPostgresUserRepository(conn)
	accountRepo := outbound.NewPostgresAccountRepository(conn)
	entryRepo := outbound.NewPostgresEntryRepository(conn)
	sessionRepo := outbound.NewPostgresSessionRepository(conn)

	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	entryService := service.NewEntryService(entryRepo)
	authService := service.NewAuthService(userRepo, sessionRepo, tokenMaker, config.ACCESS_TOKEN_DURATION, config.REFRESH_TOKEN_DURATION)

	userAdapter := httpAdapter.NewUserAdapter(*userService)
	accountAdapter := httpAdapter.NewAccountAdapter(*accountService)
	entryAdapter := httpAdapter.NewEntryAdapter(*entryService)
	authAdapter := httpAdapter.NewAuthAdapter(*authService)

	router := httpAdapter.NewRouter(tokenMaker, userAdapter, accountAdapter, entryAdapter, authAdapter)

	fmt.Println("Server starting on port " + config.PORT)
	if err := router.Run(":" + config.PORT); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
