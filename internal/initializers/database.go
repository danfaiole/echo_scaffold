package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB creates a pool to connect into a postgres database
// and returns it. Please always remember to defer close it
// after using.
func ConnectDB() *pgxpool.Pool {
	configString := fmt.Sprintf(
		"user=%v host=%v port=%v dbname=%v pool_max_conns=%v",
		os.Getenv("GO_POSTGRES_USER"),
		os.Getenv("GO_POSTGRES_HOST"),
		os.Getenv("GO_POSTGRES_PORT"),
		os.Getenv("GO_POSTGRES_DBNAME"),
		os.Getenv("GO_POSTGRES_POOL"),
	)

	dbpool, err := pgxpool.New(context.Background(), configString)

	if err != nil {
		log.Fatal("Error connecting to database, please check env vars or values to connect")
	}

	return dbpool
}
