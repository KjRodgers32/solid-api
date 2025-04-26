run:
	air

build:
	go build .

run-build:
	sudo ./solid-api
air:
	go install github.com/air-verse/air@latest
