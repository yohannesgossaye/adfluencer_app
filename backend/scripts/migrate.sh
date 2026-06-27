#!/usr/bin/env bash
# Run migrations manually (e.g. after docker compose up, or via Beekeeper alternative).
set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATION_DIR="${BACKEND_DIR}/migration"

export PGHOST="${POSTGRES_HOST:-localhost}"
export PGPORT="${POSTGRES_PORT:-5433}"
export PGUSER="${POSTGRES_USER:-adfluence}"
export PGPASSWORD="${POSTGRES_PASSWORD:-adfluence123}"
export PGDATABASE="${POSTGRES_DB:-adfluence_db}"

if ! command -v psql >/dev/null 2>&1; then
  echo "psql not found. Install PostgreSQL client or run SQL files in Beekeeper."
  exit 1
fi

echo "==> Connecting to ${PGUSER}@${PGHOST}:${PGPORT}/${PGDATABASE}"

for file in auth auth_security auth_phone_verification admin campaigns marketplace support influencer_portfolio; do
  echo "    -> ${file}.sql"
  psql -v ON_ERROR_STOP=1 -f "${MIGRATION_DIR}/${file}.sql"
done

echo "==> Done"
echo "NOTE: If types/tables already exist, use: make migrate-security"
