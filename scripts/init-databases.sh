#!/bin/sh
set -eu

cockroach_host="cockroach:26257"

sql() {
  database="$1"
  shift
  cockroach sql --insecure --host="$cockroach_host" --database="$database" "$@"
}

migrate() {
  database="$1"
  migration_dir="$2"

  echo "Migrating $database with Goose"
  GOOSE_DRIVER=postgres \
    GOOSE_DBSTRING="postgresql://root@${cockroach_host}/${database}?sslmode=disable" \
    GOOSE_MIGRATION_DIR="$migration_dir" \
    goose -no-color up
}

seed() {
  database="$1"
  seed_dir="$2"

  if [ ! -d "$seed_dir" ]; then
    return
  fi

  find "$seed_dir" -maxdepth 1 -type f -name '*.sql' -exec basename {} \; |
    sort -n |
    while IFS= read -r seed_file; do
      echo "Loading $database/$seed_file"
      sql "$database" --file="$seed_dir/$seed_file" >/dev/null
    done
}

cockroach sql --insecure --host="$cockroach_host" --execute="
  CREATE DATABASE IF NOT EXISTS school_administration;
  CREATE DATABASE IF NOT EXISTS school_user;
" >/dev/null

migrate school_administration /workspace/services/school-administration-api/database/migrations/cockroachdb
migrate school_user /workspace/services/school-user-api/database/migrations/cockroachdb

seed school_administration /workspace/services/school-administration-api/database/seeders/cockroachdb
seed school_user /workspace/services/school-user-api/database/seeders/cockroachdb

echo "Local databases are ready"
