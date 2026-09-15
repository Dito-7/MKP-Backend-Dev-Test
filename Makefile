.PHONY: run build test docker-up docker-down clean

run:
	go run cmd/api/main.go

build:
	go build -o bin/cinema-api cmd/api/main.go

test:
	go test -v ./...

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v

clean:
	rm -rf bin
