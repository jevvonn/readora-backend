run:
	@go run cmd/api/main.go

run-worker:
	@go run cmd/api/main.go

migrate-up:
	@go run cmd/migrate/main.go -m up

migrate-down:
	@go run cmd/migrate/main.go -m down

seed:
	@go run cmd/migrate/main.go -s

docs-generate:
	@swag init -g cmd/api/main.go