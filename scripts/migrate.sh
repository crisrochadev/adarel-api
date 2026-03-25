#!/usr/bin/env bash
set -euo pipefail

ACTION="${1:-up}"
MIGRATIONS_DIR="database/migrations"

: "${DB_HOST:=localhost}"
: "${DB_PORT:=5435}"
: "${DB_USER:=postgres}"
: "${DB_PASSWORD:=postgres}"
: "${DB_NAME:=cms}"

export PGPASSWORD="$DB_PASSWORD"

run_file() {
  local db="$1"
  local file="$2"
  echo "==> running ${file} on ${db}"
  psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$db" -v ON_ERROR_STOP=1 -f "$file"
}

if [[ "$ACTION" == "up" ]]; then
  run_file postgres "$MIGRATIONS_DIR/000001_create_database.up.sql" || true
  run_file "$DB_NAME" "$MIGRATIONS_DIR/000002_create_tables.up.sql"
elif [[ "$ACTION" == "down" ]]; then
  run_file "$DB_NAME" "$MIGRATIONS_DIR/000002_create_tables.down.sql"
  run_file postgres "$MIGRATIONS_DIR/000001_create_database.down.sql"
else
  echo "Usage: $0 [up|down]"
  exit 1
fi
