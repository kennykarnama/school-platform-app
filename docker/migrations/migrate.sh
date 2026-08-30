#!/bin/sh
set -eu

: "${ADMIN_DATABASE_URL:?ADMIN_DATABASE_URL is required}"
: "${USER_DATABASE_URL:?USER_DATABASE_URL is required}"

echo "Applying school administration migrations"
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="$ADMIN_DATABASE_URL" \
GOOSE_MIGRATION_DIR=/migrations/administration \
goose -no-color up

echo "Applying school user migrations"
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="$USER_DATABASE_URL" \
GOOSE_MIGRATION_DIR=/migrations/user \
goose -no-color up

echo "All database migrations are current"
