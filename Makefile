include .env

migrate-up:
	goose -dir sql/schema postgres $(DB_URL) up

migrate-down:
	goose -dir sql/schema postgres $(DB_URL) down

migrate-reset:
	goose -dir sql/schema postgres $(DB_URL) reset

build:
	go build -o out/chirpy .

run:
	go run .

test:
	go test -v -cover ./...
