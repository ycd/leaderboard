
run-dev:
	docker-compose down
	docker build . -t leaderboard
	docker-compose up

test: 
	docker build . -t leaderboard 
	docker-compose up 
	go test ./...
	clean

insert-mock-data:
	./scripts/insert_mock_data.sh $(BASE_URL)

clean-db:
	docker-compose down

clean: clean-db