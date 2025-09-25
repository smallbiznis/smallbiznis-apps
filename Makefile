# Default values for openapi-generator environment variables
GIT_USER_ID ?= taufiktriantono
GIT_REPO_ID ?= platform-sdk
GENERATOR ?= openapi-generator
GENERATOR_CMD = $(GENERATOR) generate
INPUT_DIR ?= api-specs
OUTPUT_DIR ?= gen
GENERATOR_ARGS = --git-user-id=$(GIT_USER_ID) \
                 --git-repo-id=$(GIT_REPO_ID) \
                 --additional-properties=packageName=client,enumClassPrefix=true

# Default values for golang migrate environment variables
DB_HOST ?= localhost
DB_PORT ?= 5433
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= postgres
DB_SSLMODE ?= disable

# Migration path (default to migrations directory)
MIGRATION_PATH ?= migrations

# DB Connection URL
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Migrate binary location (adjust if needed)
MIGRATE_CMD=migrate -path $(MIGRATION_PATH) -database "$(DB_URL)"

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create test-e2e bundle bundle-docs clean gen-mock $(APIS)

## Run all up migrations
migrate-up:
	$(MIGRATE_CMD) up

## Rollback all migrations
migrate-down:
	$(MIGRATE_CMD) down

## Force set migration version (useful when out of sync)
migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE_CMD) force $$version

## Print current migration version
migrate-version:
	$(MIGRATE_CMD) version

## Create new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_PATH) -format "20060102150405" $$name

gen-mock:
	go generate ./...