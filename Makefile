
run-tanx:
	@echo "Running the application"
	@go run ./cmd/tanx/main.go

run-mail:
	@echo "Running the application"
	@go run ./cmd/email/main.go

build:
	@echo "Building the application"
	@go build -o ./bin/ ./cmd/*


# Make services up
up:
	@echo "Starting the services"
	@docker compose -f ./build/compose/docker-compose.yaml up -d

# Make services down
down:
	@echo "Stopping the services"
	@docker compose -f ./build/compose/docker-compose.yaml down --volumes


# Unit tests
unit-test:
	@echo "Running unit tests..."
	@go test -v ./...
	@echo "Unit tests completed."


# Make migrations
migrate-up:
	@echo "Making migrations..."
	@migrate -path ./internal/database/mysql/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/tanx" -verbose up
	@echo "Migrations completed."

# Delete migrations
migrate-down:
	@echo "Deleting migrations..."
	@migrate -path ./internal/database/mysql/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/tanx" -verbose down

# Make migrations Devlopment environment
migrate-up-dev:
	@echo "Making migrations..."
	@migrate -path ./internal/database/mysql/migrations -database "mysql://root:password@tcp(localhost:3306)/tanx" -verbose up
	@echo "Migrations completed."

# Delete migrations Devlopment environment
migrate-down-dev:
	@echo "Deleting migrations..."
	@migrate -path ./internal/database/mysql/migrations -database "mysql://root:password@tcp(localhost:3306)/tanx" -verbose down

# SQLC generate
sqlc-gen:
	@echo "Generating SQLC..."
	@sqlc generate -f ./internal/database/mysql/sqlc/sqlc.yaml
	@echo "SQLC generation completed."