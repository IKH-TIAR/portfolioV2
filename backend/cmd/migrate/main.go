package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ikhtiar/portfolio-backend/internal/migrate"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := migrate.PoolFromURL(ctx, connString)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := migrate.Run(ctx, pool); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	fmt.Println("migrations up to date")
}
