build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

etl-load:
	docker-compose run --rm etl load_to_neo4j.py

etl-generate:
	docker-compose run --rm etl generate_csv_data.py

etl-refresh:
	make etl-generate && make etl-load

start:
	make build && make up && make etl-refresh

test:
	docker-compose run --rm api-test

clean:
	docker-compose down -v --remove-orphans
