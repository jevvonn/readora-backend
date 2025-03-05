run:
	@go run cmd/api/main.go

run-worker:
	@go run worker/server.go

migrate-up:
	@go run cmd/api/main.go -m up

migrate-down:
	@go run cmd/api/main.go -m down

seed:
	@go run cmd/api/main.go -s

docs-generate:
	@swag init -g cmd/api/main.go