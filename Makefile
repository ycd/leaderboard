
run-dev:
	docker-compose down
	docker build . -t leaderboard
	docker-compose up

run-test-env: 
	docker-compose -f dc.test.yml down
	docker-compose -f dc.test.yml up 

run-tests:
	REDIS_IP=localhost:6379 POSTGRES_USERNAME=postgres POSTGRES_PASSWORD=postgres POSTGRES_IP=localhost:5432 go test -p 1 ./... -coverprofile /dev/null

build:
	docker build . -t leaderboard

insert-mock-data:
	./scripts/insert_mock_data.sh $(BASE_URL)

clean: 
	docker-compose down

push-image: 
	docker build . -t leaderboard
	docker tag leaderboard:latest eu.gcr.io/leaderboard-312410/leaderboard:latest
	docker push eu.gcr.io/leaderboard-312410/leaderboard:latest