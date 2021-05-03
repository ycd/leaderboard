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
			SELECT COALESCE(SUM(s.point), 0)
			FROM scores s
			WHERE u.user_id=s.user_id
		) AS points
	FROM users u
	`

	CreateLeaderboardTable = `
	CREATE VIEW leaderboard AS
	SELECT 
		*,
		DENSE_RANK() OVER (ORDER BY points desc) as rank 
	FROM UsersWithScores`

	GetLeaderboard = `
	SELECT 
		rank,
		points,
		name as display_name,
		country
	FROM leaderboard
	ORDER BY rank`

	GetLeaderboardWithCountry = `
	SELECT 
		DENSE_RANK() OVER (ORDER BY points desc) as rank,
		points,
		name as display_name,
		country
	FROM UsersWithScores
	WHERE country = $1
	ORDER BY rank
	`

	InsertScore = `
	INSERT INTO scores
	VALUES ($1, $2, $3)`

	NewUser = `
	INSERT INTO users
	VALUES ($1, $2, $3)`

	GetUser = `
	SELECT 
		user_id,
		name as display_name,
		points,
		rank
	FROM leaderboard
	WHERE user_id = $1`
)
