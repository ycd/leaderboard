package queries

var (
	CreateScoresTable = `
	CREATE TABLE IF NOT EXISTS scores (
		user_id VARCHAR(40) UNIQUE NOT NULL,
		point SMALLINT,
		timestamp TIMESTAMP NOT NULL
	)
	`

	CreateUserTable = `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(40) UNIQUE NOT NULL,
		name VARCHAR(40) UNIQUE NOT NULL,
		created_at TIMESTAMP,
		country VARCHAR(10)
	)
	`

	CreateUserWithScoresView = `
	CREATE VIEW UsersWithScores AS
	SELECT *, (SELECT SUM(point)+1000 as point FROM scores s WHERE u.id=s.user_id ) as point FROM users u
	`

	CreateLeaderboardTable = `
	CREATE MATERIALIZED VIEW leaderboard AS
	SELECT 
		*,
		RANK() OVER (ORDER BY point desc) as rank 
	FROM UsersWithScores
	`

	GetLeaderboard = `
	SELECT *, 
	FROM leaderboard`

	GetLeaderboardWithCountry = `
	SELECT *, 
	FROM leaderboard
	WHERE COUNTRY = %s`

	SubmitScore = `
	INSERT INTO scores
	VALUES ()
	`

	// TODO

)
