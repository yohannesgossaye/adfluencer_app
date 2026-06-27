#!/usr/bin/env bash
# Apply influencer_portfolio.sql only (existing DB).
set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_FILE="${BACKEND_DIR}/migration/influencer_portfolio.sql"

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
echo "    -> influencer_portfolio.sql"
psql -v ON_ERROR_STOP=1 -f "$SQL_FILE"
echo "==> Done (influencer_portfolio_items)"
