# Use the .envrc file
include .envrc

########################################################################################################
# Commands to run the application and tests
########################################################################################################
# run/api: run the cmd/api application
.PHONY: run/api
run/api:
        @echo "Starting API server on port $(PORT) in $(ENV) mode..."
        @go run ./cmd/api \
        -port=$(PORT) \
        -env=$(ENV) \
        -db-dsn=$(DB_DSN) \
        -db-max-open-conns=$(DB_MAX_OPEN_CONNS) \
        -db-max-idle-conns=$(DB_MAX_IDLE_CONNS) \
        -db-max-idle-time=$(DB_MAX_IDLE_TIME) \
        -cors-trusted-origins="$(CORS_TRUSTED_ORIGINS)\
        -rate-limiter-enabled=$(RATE_LIMITER_ENABLED) \
        -rate-limiter-rps=$(RATE_LIMITER_RPS) \
        -rate-limiter-burst=$(RATE_LIMITER_BURST)"

## run/tests: run the tests
.PHONY: run/tests
run/tests:
        @echo "Running tests..."
        @go test ./...



########################################################################################################
# PostgreSQL commands
########################################################################################################
# Login to psql
.PHONY: psql/login
psql/login:
		psql "$(DB_DSN)"

# login to postgresql as sudo
.PHONY: psql/sudo
psql/sudo:
        sudo -u postgres psql

########################################################################################################
# Database migration commands
########################################################################################################
# Create a new migration file
.PHONY: migration/create
migration/create:
        @if [ -z "$(name)" ]; then \
                echo "Error: Please provide a name for the migration using 'make migration/create name=your_migration_name'"; \
                exit 1; \
        fi
        @if [ ! -d "./migrations" ]; then mkdir ./migrations; fi
        migrate create -seq -ext=.sql -dir=./migrations $(name)

# Apply all up migrations
.PHONY: migration/up
migration/up:
        migrate -path ./migrations -database "$(DB_DSN)" up 1

# Apply all down 1 migrations
.PHONY: migration/down
migration/down:
        migrate -path ./migrations -database "$(DB_DSN)" down 1

# fix and reapply the last migration and fix dirty state
.PHONY: migration/fix
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
#                 echo "Running up migration..."; \
#               migrate -path ./migrations -database "${DB_DSN}" up; \
        else \
                echo "No dirty migration found"; \
        fi
        @rm -f /tmp/migrate_version

