package storage

import (
	"database/sql"
	"log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ycd/leaderboard/internal/storage/queries"
)

func NewClickHouseConnection() (*sql.DB, error) {
	options := getClickHouseOptions()

	db := clickhouse.OpenDB(options)
	if os.Getenv("REFRESH") == "true" {
		dropQueries := []string{
			queries.DropClickhouseTableEventParticipant,
			queries.DropClickhouseTableUser,
			queries.DropClickhouseTableEvent,
			queries.DropEventLeaderboardMV,
			queries.DropEventCountryLeaderboardMV,
		}
		if err := executeQueries(db, dropQueries); err != nil {
			log.Println("Error dropping ClickHouse tables:", err)
		}

		createQueries := []string{
			queries.CreateEventParticipantTable,
			queries.CreateUserTable,
			queries.CreateEventTable,
			queries.CreateEventLeaderboardMV,
		}
		if err := executeQueries(db, createQueries); err != nil {
			return nil, err
		}

		log.Println("Recreated ClickHouse tables based on CREATE queries.")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getClickHouseOptions() *clickhouse.Options {
	return &clickhouse.Options{
		Addr: []string{os.Getenv("CLICKHOUSE_ADDR")},
		Auth: clickhouse.Auth{
			Database: os.Getenv("CLICKHOUSE_DB"),
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	}
}

func executeQueries(db *sql.DB, queries []string) error {
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error executing query: %v | Query: %s\n", err, query)
			// Continue executing other queries even if one fails
			continue
		}
		log.Printf("Executed query successfully: %s\n", query)
	}
	return nil
}
