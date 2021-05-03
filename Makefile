
run-dev:
	docker-compose down
	docker build . -t leaderboard
	docker-compose up

test-db: 
	docker build . -t leaderboard
	docker-compose -f dc.test.yml up 

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