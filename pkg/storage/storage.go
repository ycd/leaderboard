package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ycd/leaderboard/pkg/queries"
)

// Storage holds the connection pool between PostgreSQL.
// With a bunch of methods on top of it.
type Storage struct {
	connection *pgxpool.Pool
	cache      *redis.Client
}

func init() {
	ctx := context.Background()
	storage := NewStorage(ctx)
	storage.CreateTables(ctx)
}

// NewStorage creates a new Storage
func NewStorage(ctx context.Context) *Storage {
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresIPAddr := os.Getenv("POSTGRES_IP")
	redisIPAddr := os.Getenv("REDIS_IP")

	conn, err := pgxpool.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s/postgres", postgresUsername, postgresPassword, postgresIPAddr))
	if err != nil {
		log.Fatalf("ERROR: unable to connect PSQL: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisIPAddr,
		Password: "",
		DB:       0,
	})

	return &Storage{
		connection: conn,
		cache:      rdb,
	}
}

// Conn exposes the underlying connection for health checking purposes.
func (s *Storage) Conn() *pgxpool.Pool {
	return s.connection
}

func (s *Storage) Cache() *redis.Client {
	return s.cache
}

// Create the tables on startup, this function intented to run only on startup.
func (s *Storage) CreateTables(ctx context.Context) {
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
