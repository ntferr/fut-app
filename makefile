setup: migrate
	go run cmd/migrate.go

up-db:
	docker-compose -f infra/docker-compose.yaml up -d postgres omnidb