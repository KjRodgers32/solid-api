include .env
export

DB_URL=postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB}?sslmode=disable
run:
	air

build:
	go build .

run-build:
	sudo ./solid-api

migrate-up:
	migrate -path db/migrations -database "${DB_URL}" up

migrate-down:
	migrate -path db/migrations -database "${DB_URL}" down