include .envrc

.PHONY: run/api run/tests run/api/win psql/login psql/sudo migration/create migration/up migration/down migration/fix db/migrations/up
run/api:
	@echo "Starting API server on port $(PORT) in $(ENV) mode..."
	@go run ./cmd/api \
		-port $(PORT) \
		-env $(ENVIRONMENT) \
		-db-dsn "$(DB_DSN)" \
		-db-max-open-conns $(DB_MAX_OPEN_CONNS) \
		-db-max-idle-conns $(DB_MAX_IDLE_CONNS) \
		-db-max-idle-time $(DB_MAX_IDLE_TIME) \
		-cors-trusted-origins "$(CORS_ALLOWED_ORIGINS)" \
		-limiter-enabled=$(RATE_LIMITER_ENABLED) \
		-limiter-rps $(RATE_LIMITER_RPS) \
		-limiter-burst $(RATE_LIMITER_BURST) \
		-smtp-host "$(SMTP_HOST)" \
		-smtp-port $(SMTP_PORT) \
		-smtp-username "$(SMTP_USERNAME)" \
		-smtp-password "$(SMTP_PASSWORD)" \
		-smtp-sender "$(SMTP_SENDER)"

run/api/win:
	@powershell -Command "Get-Content .envrc | ForEach-Object { if ($$_ -match '^([^=]+)=(.*)$$') { $$value = $$matches[2] -replace '^\"(.*)\"$$', '$$1'; [System.Environment]::SetEnvironmentVariable($$matches[1], $$value, 'Process') } }; go run ./cmd/api"

run/tests:
	@echo "Running tests..."
	@go test ./...

db/migrations/up:
	@migrate -path ./migrations -database "$(DB_DSN)" up

psql/login:
	psql "$(DB_DSN)"

psql/sudo:
	sudo -u postgres psql

migration/create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a name for the migration using 'make migration/create name=your_migration_name'"; \
		exit 1; \
	fi
	@if [ ! -d "./migrations" ]; then mkdir ./migrations; fi
	migrate create -seq -ext=.sql -dir=./migrations $(name)

migration/up:
	migrate -path ./migrations -database "$(DB_DSN)" up 

migration/down:
	migrate -path ./migrations -database "$(DB_DSN)" down 

migration/fix:
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