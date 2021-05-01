package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ycd/leaderboard/pkg/queries"
)

// Storage holds the connection pool between PostgreSQL.
// With a bunch of methods on top of it.
type Storage struct {
	connection *pgxpool.Pool
}

func init() {
	ctx := context.Background()
	storage := NewStorage(ctx)
	storage.createTables(ctx)
}

// NewStorage creates a new Storage
func NewStorage(ctx context.Context) *Storage {
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresDBName := os.Getenv("POSTGRES_DB")
	postgresIPAddr := os.Getenv("POSTGRES_IP")

	conn, err := pgxpool.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s", postgresUsername, postgresPassword, postgresIPAddr, postgresDBName))
	if err != nil {
		log.Fatalf("ERROR: unable to connect PSQL: %v", err)
	}

	return &Storage{
		connection: conn,
	}
}

// Create the tables on startup, this function intented to run only on startup.
func (s *Storage) createTables(ctx context.Context) {
	for _, query := range []string{
		queries.CreateScoresTable,
		queries.CreateUserTable,
		queries.CreateUserWithScoresView,
		queries.CreateLeaderboardTable,
	} {
		rows, err := s.connection.Query(ctx, query)
		if err != nil {
			log.Printf("Failed to create: %v", err)
		}
		log.Printf("Created table successfully: %v", rows)
	}
}
