
run-dev:
	docker-compose down
	docker build . -t leaderboard
	docker-compose up

test: 
	docker build . -t leaderboard
	docker-compose up 
	go test ./...
	docker-compose down

build:
	docker build . -t leaderboard

insert-mock-data:
	./scripts/insert_mock_data.sh $(BASE_URL)

clean: 
	docker-compose down