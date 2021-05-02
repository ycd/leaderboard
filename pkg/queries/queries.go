package queries

var (
	CreateScoresTable = `
	CREATE TABLE IF NOT EXISTS scores (
		user_id VARCHAR(40) NOT NULL,
		point SMALLINT,
		timestamp INTEGER NOT NULL
	)`

	CreateUserTable = `
	CREATE TABLE IF NOT EXISTS users (
		user_id VARCHAR(40) UNIQUE NOT NULL,
		name VARCHAR(40) UNIQUE NOT NULL,
		country VARCHAR(10)
	)`

	CreateUserWithScoresView = `
	CREATE VIEW UsersWithScores AS
	SELECT 
		u.user_id,
		u.name,
		u.country,
		(
			SELECT SUM(s.point)
			FROM scores s
			WHERE u.user_id=s.user_id
		) AS point
	FROM users u
	`

	CreateLeaderboardTable = `
	CREATE VIEW leaderboard AS
	SELECT 
		*,
		RANK() OVER (ORDER BY point desc) as rank 
	FROM UsersWithScores`

	GetLeaderboard = `
	SELECT *
	FROM leaderboard`

	GetLeaderboardWithCountry = `
	SELECT *
	FROM leaderboard
	WHERE COUNTRY = $1`

	InsertScore = `
	INSERT INTO scores
	VALUES ($1, $2, $3)`

	NewUser = `
	INSERT INTO users
	VALUES ($1, $2, $3)`

	GetUser = `
	SELECT * 
	FROM leaderboard
	WHERE user_id=$1`
)
