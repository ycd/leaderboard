package queries

const (
	CreateEventParticipantTable = `
    CREATE TABLE IF NOT EXISTS event_participants (
        event_id String,
        user_id String,
        score Int32,
        country String,
        timestamp DateTime DEFAULT now()
    ) ENGINE = MergeTree()
    ORDER BY (event_id, user_id, timestamp)
    `
)

const (
	CreateUserTable = `
		CREATE TABLE IF NOT EXISTS users (
			id String,
			username String,
			country String,
			level Int32,
			coin Int32,
			is_banned UInt8 DEFAULT 0
		) ENGINE = MergeTree()
		ORDER BY id;
		`
)

const (
	CreateEventTable = `
		CREATE TABLE IF NOT EXISTS events (
			id String,
			name String,
			start_time DateTime,
			end_time DateTime
		) ENGINE = MergeTree()
		ORDER BY id;
		`
)

const (
	CreateEventLeaderboardMV = `
    CREATE MATERIALIZED VIEW IF NOT EXISTS event_leaderboard_mv
    ENGINE = MergeTree()
    ORDER BY (event_id, total_score, user_id)
    AS SELECT
        event_id,
        user_id,
        SUM(score) AS total_score,
        any(country) AS country
    FROM event_participants
    GROUP BY event_id, user_id
	ORDER BY total_score DESC
		`
)

const (
	CreateSqliteUserTable = `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		country TEXT NOT NULL,
		level INTEGER DEFAULT 0,
		coin INTEGER DEFAULT 0,
		is_banned INTEGER DEFAULT 0
	);
	`
)

const (
	CreateSqliteRewardsTable = `
    CREATE TABLE IF NOT EXISTS rewards (
        event_id TEXT,
        user_id TEXT,
        claimed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        reward_amount INTEGER,
        PRIMARY KEY (event_id, user_id)
    )
	`
)

const (
	DropClickhouseTableEventParticipant = `
	DROP TABLE IF EXISTS event_participants;
	`
	DropClickhouseTableUser = `
	DROP TABLE IF EXISTS users;
	`
	DropClickhouseTableEvent = `
	DROP TABLE IF EXISTS events;
	`
	// Drop statements for Materialized Views
	DropEventLeaderboardMV = `
	DROP TABLE IF EXISTS event_leaderboard_mv;
	`
	DropEventCountryLeaderboardMV = `
	DROP TABLE IF EXISTS event_country_leaderboard_mv;
	`
)

const (
	DropSqliteTables = `
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS rewards;
	`
)
