#!/usr/bin/env bash
# Apply auth_security.sql only (existing DB — does not re-run auth.sql).
set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_FILE="${BACKEND_DIR}/migration/auth_security.sql"

export PGHOST="${POSTGRES_HOST:-localhost}"
export PGPORT="${POSTGRES_PORT:-5433}"
export PGUSER="${POSTGRES_USER:-adfluence}"
export PGPASSWORD="${POSTGRES_PASSWORD:-adfluence123}"
export PGDATABASE="${POSTGRES_DB:-adfluence_db}"

if ! command -v psql >/dev/null 2>&1; then
  echo "psql not found. Install PostgreSQL client."
  exit 1
fi

echo "==> Connecting to ${PGUSER}@${PGHOST}:${PGPORT}/${PGDATABASE}"
echo "    -> auth_security.sql"
psql -v ON_ERROR_STOP=1 -f "$SQL_FILE"
echo "==> Done (auth_login_lockouts, auth_rate_limit_events)"
