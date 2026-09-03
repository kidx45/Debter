package initiator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	httpAdapter "github.com/kidx45/Debter/internal/adapter/inbound/http"
	outboundAdapter "github.com/kidx45/Debter/internal/adapter/outbound"
	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/service"
	"github.com/kidx45/Debter/internal/util/token"
)

type Server struct {
	config  domain.Config
	db      *sql.DB
	httpSrv *http.Server
	router  *gin.Engine
}

func NewServer(config domain.Config) (*Server, error) {
	conn, err := sql.Open(config.DB_DRIVER_NAME, config.DB_URL)
	if err != nil {
		return nil, fmt.Errorf("error setting up db connection: %w", err)
	}

	queries := db.New(conn)

	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		if pqErr := conn.Close(); pqErr != nil {
			return nil, fmt.Errorf("error setting up token maker: %w, error closing the db :%w", err, pqErr)
		}
		return nil, fmt.Errorf("error setting up token maker: %w", err)
	}

	userRepo := outboundAdapter.NewPostgresUserRepository(queries)
	accountRepo := outboundAdapter.NewPostgresAccountRepository(queries)
	entryRepo := outboundAdapter.NewPostgresEntryRepository(queries, conn)
	sessionRepo := outboundAdapter.NewPostgresSessionRepository(queries)

	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	entryService := service.NewEntryService(entryRepo)
	authService := service.NewAuthService(userRepo, sessionRepo, tokenMaker, config.ACCESS_TOKEN_DURATION, config.REFRESH_TOKEN_DURATION)

	userAdapter := httpAdapter.NewUserAdapter(*userService)
	accountAdapter := httpAdapter.NewAccountAdapter(*accountService)
	entryAdapter := httpAdapter.NewEntryAdapter(*entryService)
	authAdapter := httpAdapter.NewAuthAdapter(*authService)
	healthAdapter := httpAdapter.NewHealthAdapter(conn.PingContext)

	router := httpAdapter.NewRouter(tokenMaker, userAdapter, accountAdapter, entryAdapter, authAdapter, healthAdapter)

	httpSrv := &http.Server{
		Addr:    ":" + config.PORT,
		Handler: router,
	}

	return &Server{
		config:  config,
		db:      conn,
		httpSrv: httpSrv,
		router:  router,
	}, nil
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server starting on port %s\n", s.config.PORT)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal %s, shutting down gracefully...\n", sig)
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	fmt.Println("Server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %s", err)
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db close error: %w", err)
	}
	return nil
}
