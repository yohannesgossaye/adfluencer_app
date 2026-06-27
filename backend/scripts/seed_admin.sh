#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_FILE="${ROOT}/migration/seed_admin.sql"

PGHOST="${POSTGRES_HOST:-localhost}"
PGPORT="${POSTGRES_PORT:-5433}"
PGUSER="${POSTGRES_USER:-adfluence}"
PGPASSWORD="${POSTGRES_PASSWORD:-adfluence123}"
PGDATABASE="${POSTGRES_DB:-adfluence_db}"

export PGPASSWORD

echo "Seeding admin user (admin@adfluence.com)..."
psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -f "$SQL_FILE"
echo "Done."
