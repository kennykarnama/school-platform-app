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

cockroach sql --insecure --host="$cockroach_host" --execute="
  CREATE DATABASE IF NOT EXISTS school_administration;
  CREATE DATABASE IF NOT EXISTS school_user;
" >/dev/null

migrate school_administration /workspace/services/school-administration-api/database/migrations/cockroachdb
migrate school_user /workspace/services/school-user-api/database/migrations/cockroachdb

echo "Loading administration seed data"
sql school_administration --file=/workspace/services/school-administration-api/database/seeders/cockroachdb/1_seed-academic_year.sql >/dev/null

echo "Local databases are ready"
