package storage

import (
	"database/sql"
	"os"

	"github.com/ycd/leaderboard/internal/storage/queries"
	_ "modernc.org/sqlite"
)

func NewSQLiteConnection() (*sql.DB, error) {
	dsn := getSQLiteDataSourceName()
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if os.Getenv("REFRESH") == "true" {
		queries := []string{
			queries.DropSqliteTables,
		}
		if err := executeQueries(db, queries); err != nil {
			return nil, err
		}
	}

	queries := []string{
		queries.CreateSqliteUserTable,
		queries.CreateSqliteRewardsTable,
	}

	if err := executeQueries(db, queries); err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func getSQLiteDataSourceName() string {
	return "file:leaderboard.db?cache=shared&_journal_mode=WAL&synchronous=OFF"
}
