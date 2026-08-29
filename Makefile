SHELL := /bin/sh

.PHONY: setup deps-up deps-down db-reset run-administration run-user run-ui dev test build

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

