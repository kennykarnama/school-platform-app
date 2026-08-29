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

  sql "$database" --execute="
    CREATE TABLE IF NOT EXISTS schema_migrations (
      version STRING PRIMARY KEY,
      applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
  " >/dev/null

  find "$migration_dir" -maxdepth 1 -type f -name '*.sql' -exec basename {} \; |
    sort -n |
    while IFS= read -r migration; do
      applied="$(sql "$database" --format=csv --execute="SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = '$migration');" | tail -n 1 | tr -d '\r')"
      if [ "$applied" = "t" ]; then
        echo "Skipping $database/$migration (already applied)"
        continue
      fi

      echo "Applying $database/$migration"
      sed -e '$a\' "$migration_dir/$migration" | sql "$database" >/dev/null
      sql "$database" --execute="INSERT INTO schema_migrations (version) VALUES ('$migration');" >/dev/null
    done
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
