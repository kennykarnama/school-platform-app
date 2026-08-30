SHELL := /bin/sh

.PHONY: setup deps-up deps-down db-reset db-migrate-administration db-migrate-user db-status-administration db-status-user db-validate db-new-administration db-new-user run-administration run-user run-ui dev test build

setup:
	npm --prefix apps/school-administration-ui ci

deps-up:
	docker compose up -d --wait cockroach redis
	docker compose run --rm db-init

deps-down:
	docker compose down --remove-orphans

db-reset:
	docker compose down --volumes --remove-orphans
	docker compose up -d --wait cockroach redis
	docker compose run --rm db-init

db-migrate-administration:
	@test -n "$$DATABASE_URL" || (echo "DATABASE_URL is required" && exit 1)
	@GOOSE_DBSTRING="$$DATABASE_URL" docker compose --profile tools run --rm -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING goose -no-color -dir /workspace/services/school-administration-api/database/migrations/cockroachdb up

db-migrate-user:
	@test -n "$$DATABASE_URL" || (echo "DATABASE_URL is required" && exit 1)
	@GOOSE_DBSTRING="$$DATABASE_URL" docker compose --profile tools run --rm -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING goose -no-color -dir /workspace/services/school-user-api/database/migrations/cockroachdb up

db-status-administration:
	@test -n "$$DATABASE_URL" || (echo "DATABASE_URL is required" && exit 1)
	@GOOSE_DBSTRING="$$DATABASE_URL" docker compose --profile tools run --rm -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING goose -no-color -dir /workspace/services/school-administration-api/database/migrations/cockroachdb status

db-status-user:
	@test -n "$$DATABASE_URL" || (echo "DATABASE_URL is required" && exit 1)
	@GOOSE_DBSTRING="$$DATABASE_URL" docker compose --profile tools run --rm -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING goose -no-color -dir /workspace/services/school-user-api/database/migrations/cockroachdb status

db-validate:
	docker compose --profile tools run --rm --no-deps goose -no-color -dir /workspace/services/school-administration-api/database/migrations/cockroachdb validate
	docker compose --profile tools run --rm --no-deps goose -no-color -dir /workspace/services/school-user-api/database/migrations/cockroachdb validate

db-new-administration:
	@test -n "$(NAME)" || (echo "NAME is required" && exit 1)
	docker compose --profile tools run --rm --no-deps goose -s -dir /workspace/services/school-administration-api/database/migrations/cockroachdb create "$(NAME)" sql

db-new-user:
	@test -n "$(NAME)" || (echo "NAME is required" && exit 1)
	docker compose --profile tools run --rm --no-deps goose -s -dir /workspace/services/school-user-api/database/migrations/cockroachdb create "$(NAME)" sql

run-administration:
	set -a; . ./config/local/administration.env; set +a; cd services/school-administration-api && go run .

run-user:
	set -a; . ./config/local/user.env; set +a; cd services/school-user-api && go run .

run-ui:
	cd apps/school-administration-ui && SCHOOL_ADMINISTRATION_API_BASE_URL=http://localhost:8081 npm run start-dev

dev: deps-up
	./scripts/dev.sh

test:
	cd services/school-administration-api && go test ./...
	cd services/school-user-api && go test ./...
	npm --prefix apps/school-administration-ui test

build:
	cd services/school-administration-api && go build ./...
	cd services/school-user-api && go build ./...
	SCHOOL_ADMINISTRATION_API_BASE_URL=http://localhost:8081 npm --prefix apps/school-administration-ui run build
