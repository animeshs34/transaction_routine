.PHONY: build run test docker-up docker-down

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./... -v

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
