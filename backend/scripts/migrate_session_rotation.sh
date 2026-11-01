#!/usr/bin/env bash
# Apply session_rotation.sql only (existing DB). Safe to re-run.
# Local:  ./scripts/migrate_session_rotation.sh
# Remote: DATABASE_URL="postgres://..." ./scripts/migrate_session_rotation.sh
set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_FILE="${BACKEND_DIR}/migration/session_rotation.sql"

if ! command -v psql >/dev/null 2>&1; then
  echo "psql not found. Install PostgreSQL client."
  exit 1
fi

if [ -n "${DATABASE_URL:-}" ]; then
  echo "==> Applying session_rotation.sql to \$DATABASE_URL"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$SQL_FILE"
else
  export PGHOST="${POSTGRES_HOST:-localhost}"
  export PGPORT="${POSTGRES_PORT:-5433}"
  export PGUSER="${POSTGRES_USER:-adfluence}"
  export PGPASSWORD="${POSTGRES_PASSWORD:-adfluence123}"
  export PGDATABASE="${POSTGRES_DB:-adfluence_db}"
  echo "==> Connecting to ${PGUSER}@${PGHOST}:${PGPORT}/${PGDATABASE}"
  psql -v ON_ERROR_STOP=1 -f "$SQL_FILE"
fi
echo "==> Done (user_sessions rotation columns + refresh hash index)"
