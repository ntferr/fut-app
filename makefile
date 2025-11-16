
up-database:
	docker-compose -f infra/docker-compose.yaml up -d postgres omnidb

run-migrate:
	go run cmd/migrator/main.go