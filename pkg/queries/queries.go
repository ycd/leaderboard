package queries

var (
	CreateScoresTable = `
	CREATE TABLE IF NOT EXISTS scores (
		user_id VARCHAR(40) NOT NULL,
		point SMALLINT,
		timestamp INTEGER NOT NULL
	)
	`

	CreateUserTable = `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(40) UNIQUE NOT NULL,
		name VARCHAR(40) UNIQUE NOT NULL,
		created_at INTEGER,
		country VARCHAR(10)
	)
	`

	CreateUserWithScoresView = `
	CREATE VIEW UsersWithScores AS
	SELECT *, (SELECT SUM(point) as point FROM scores s WHERE u.id=s.user_id ) as point FROM users u
	`

	CreateLeaderboardTable = `
	CREATE MATERIALIZED VIEW leaderboard AS
	SELECT 
		*,
		RANK() OVER (ORDER BY point desc) as rank 
	FROM UsersWithScores
	`

	GetLeaderboard = `
	SELECT *
	FROM leaderboard`

	GetLeaderboardWithCountry = `
	SELECT *
	FROM leaderboard
	WHERE COUNTRY = $1`

	InsertScore = `
	INSERT INTO scores
	VALUES ($1, $2, $3)
	`

	// TODO

)
