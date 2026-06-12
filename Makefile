.PHONY: run test sqlc migrate

run:
	go run ./cmd/server

test:
	go test ./...

sqlc:
	sqlc generate

migrate:
	psql -U $${DB_USER:-postgres} -d $${DB_NAME:-userdb} -f db/migrations/001_create_users.sql
