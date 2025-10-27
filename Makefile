include .envrc

.PHONY: run/api run/tests run/api/win psql/login psql/sudo migrate/create migrate/up migrate/down migrate/fix d b/migrations/up swagger/docs pg_dump/schema run/users run/officers run/tests/internal/data run/tests/cmd/api run/populate_data
run/api:
	@echo "Starting API server on port $(PORT) in $(ENVIRONMENT) mode..."
	@go run ./cmd/api \
		-port=$(PORT) \
		-env=$(ENVIRONMENT) \
		-db-dsn=$(DB_DSN) \
		-db-max-open-conns=$(DB_MAX_OPEN_CONNS) \
		-db-max-idle-conns=$(DB_MAX_IDLE_CONNS) \
		-db-max-idle-time=$(DB_MAX_IDLE_TIME) \
		-cors-trusted-origins=$(CORS_ALLOWED_ORIGINS) \
		-limiter-enabled=$(RATE_LIMITER_ENABLED) \
		-limiter-rps=$(RATE_LIMITER_RPS) \
		-limiter-burst=$(RATE_LIMITER_BURST) \
		-smtp-host=$(SMTP_HOST) \
		-smtp-port=$(SMTP_PORT) \
		-smtp-username=$(SMTP_USERNAME) \
		-smtp-password=$(SMTP_PASSWORD) \
		-smtp-sender=$(SMTP_SENDER)

run/api/win:
	@powershell -Command "Get-Content .envrc | ForEach-Object { if ($$_ -match '^([^=]+)=(.*)$$') { $$value = $$matches[2] -replace '^\"(.*)\"$$', '$$1'; [System.Environment]::SetEnvironmentVariable($$matches[1], $$value, 'Process') } }; go run ./cmd/api"

run/tests:
	@echo "Running tests..."
	@go test ./...

db/migrations/up:
	@migrate -path ./migrations -database "$(DB_DSN)" up

psql/login/test:
	psql "$(TEST_DB_DSN)"

psql/login:
	psql "$(DB_DSN)"

psql/sudo:
	sudo -u postgres psql

migrate/create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a name for the migration using 'make migrate/create name=your_migration_name'"; \
		exit 1; \
	fi
	@if [ ! -d "./migrations" ]; then mkdir ./migrations; fi
	migrate create -seq -ext=.sql -dir=./migrations $(name)

migrate/up:
	migrate -path ./migrations -database "$(DB_DSN)" up 

migrate/down:
	migrate -path ./migrations -database "$(DB_DSN)" down

migrate/up1:
	migrate -path ./migrations -database "$(DB_DSN)" up 1

migrate/down1:
	migrate -path ./migrations -database "$(DB_DSN)" down 1

migrate/up/test:
	migrate -path ./migrations -database "$(TEST_DB_DSN)" up
migrate/down/test:
	migrate -path ./migrations -database "$(TEST_DB_DSN)" down

migrate/reset/test:
	migrate -path ./migrations -database "$(TEST_DB_DSN)" down
	migrate -path ./migrations -database "$(TEST_DB_DSN)" up

migrate/fix:
	@echo 'Checking migration status...'
	@migrate -path ./migrations -database "${DB_DSN}" version > /tmp/migrate_version 2>&1
	@cat /tmp/migrate_version
	@if grep -q "dirty" /tmp/migrate_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path ./migrations -database "${DB_DSN}" force $$version; \
		echo "Running down migration..."; \
		migrate -path ./migrations -database "${DB_DSN}" down 1; \
		echo "Running up migration..."; \
		migrate -path ./migrations -database "${DB_DSN}" up 1; \
	else \
		echo "No dirty migration found"; \
	fi
	@rm -f /tmp/migrate_version

swagger/docs:
	@echo "Generating Swagger documentation..."
	@swag init -g ./cmd/api/main.go


migrate/fix/test:
	@echo 'Checking migration status...'
	@migrate -path ./migrations -database "${TEST_DB_DSN}" version > /tmp/migrate_version 2>&1
	@cat /tmp/migrate_version
	@if grep -q "dirty" /tmp/migrate_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path ./migrations -database "${TEST_DB_DSN}" force $$version; \
		echo "Running down migration..."; \
		migrate -path ./migrations -database "${TEST_DB_DSN}" down 1; \
		echo "Running up migration..."; \
		migrate -path ./migrations -database "${TEST_DB_DSN}" up 1; \
	else \
		echo "No dirty migration found"; \
	fi
	@rm -f /tmp/migrate_version

pg_dump/schema:
	@echo "Exporting database schema to schema.sql..."
	@pg_dump -h localhost -p 5432 -d police_training -U postgres -s -F p -E UTF-8 -f ~/Projects/police_training/schema.sql

run/users:
	@echo "Running Users Population Script..."
	@go run ./cmd/populate_users/main.go

run/officers:
	@echo "Running Officers Population Script..."
	@go run ./cmd/populate_officers/main.go


run/tests:
	@echo "Running tests..."
	@go test -v ./...

run/tests/internal/data:
	@echo "Running tests for internal/data package..."
	@go test -v ./internal/data/...

run/tests/cmd/api:
	@echo "Running tests for cmd/api package..."
	@go test -v ./cmd/api/...
# FOR TESTING INDIVIDUAL PACKAGES
# go test -v ./cmd/api/
# go test -v ./internal/data/
# go test -v ./...

run/populate_data:
	@echo "Running Data Population Script..."
	@go run ./cmd/populate_data/main.go -dsn "$(TEST_DB_DSN)"

# Test API handlers only
.PHONY: test-api
test-api:
	@echo "Running API handler tests only..."
	go test -v ./cmd/api/... | tee api_test_results.txt

# Test API handlers with coverage
.PHONY: test-api-coverage
test-api-coverage:
	@echo "Running API handler tests with coverage..."
	go test -v -coverprofile=api_coverage.out ./cmd/api/... | tee api_test_results.txt
	go tool cover -html=api_coverage.out -o api_coverage.html
	@echo "Coverage report saved to api_coverage.html"

# Clean test artifacts
.PHONY: clean-test
clean-test:
	rm -f api_test_results.txt api_coverage.out api_coverage.html