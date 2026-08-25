package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/kidx45/Debter/internal/initiator"
	"github.com/kidx45/Debter/internal/util"
)

func main() {
	config, err := util.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading the Env Configuration: %s", err)
	}

	srv, err := initiator.NewServer(config)
	if err != nil {
		log.Fatalf("Error creating server: %s", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
