.PHONY: up down rebuild logs clean localbuild postgres lint

up:
	docker compose up -d

postgres:
	docker compose up -d postgres

down:
	docker compose down

rebuild:
	docker compose down -v && docker compose up -d --build && docker compose logs -f

logs:
	docker compose logs -f

clean:
	docker compose down -v

localbuild:
	go build -C cmd/main -o ../../goprocess .

lint:
	cd cmd/main && golangci-lint run ./...
