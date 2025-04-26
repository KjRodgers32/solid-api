run:
	docker-compose up --build

stop:
	docker-compose down

build:
	go build .

run-build:
	sudo ./solid-api
air:
	go install github.com/air-verse/air@latest
