package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ycd/leaderboard/pkg/queries"
)

type Storage struct {
	connection *pgxpool.Pool
}

func NewStorage(ctx context.Context) *Storage {
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresDBName := os.Getenv("POSTGRES_DB")
	postgresIPAddr := os.Getenv("POSTGRES_IP")

	conn, err := pgxpool.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s", postgresUsername, postgresPassword, postgresIPAddr, postgresDBName))
	if err != nil {
		log.Fatalf("ERROR: unable to connect PSQL: %v", err)
	}

	storage := &Storage{
		connection: conn,
	}

	storage.createTables(ctx)

	return storage
}

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

func (s *Storage) GetLeaderboard() error {
	rows, err := s.connection.Query(context.Background(), queries.GetLeaderboard)
	if err != nil {
		return err
	}

	for rows.Next() {
		v, err := rows.Values()
		if err != nil {
			log.Println("Got error:", err)
		}

		fmt.Println("Got value:", v)
	}

	return nil
}
