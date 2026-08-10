package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/kidx45/Debter/internal/adapter/outbound"
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

	_ = outbound.NewRepository(conn)
	fmt.Println("Connection Successful")
}
