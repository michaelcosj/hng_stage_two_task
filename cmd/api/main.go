package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/michaelcosj/hng-task-two/internal/server"
)

func main() {
	log.Printf("Initialising database connection\n")

	// TODO: find out why pgxpool seems to be slower than pgx
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatalf("Error initialising database: %v", err)
	}
	defer conn.Close()

	log.Printf("Starting server...\n")
	if err := server.New(conn).ListenAndServe(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
