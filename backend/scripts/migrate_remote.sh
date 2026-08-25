#!/usr/bin/env bash
# Migrate any remote PostgreSQL database (Render, Supabase, Neon, Railway, etc.)
# Usage:
#   ./scripts/migrate_remote.sh "postgres://user:password@host:5432/dbname?sslmode=require"
# Or:
#   DATABASE_URL="postgres://..." ./scripts/migrate_remote.sh

set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATION_FILE="${BACKEND_DIR}/migration/deploy_all.sql"

TARGET_URL="${1:-${DATABASE_URL:-}}"

if [ -z "$TARGET_URL" ]; then
  echo "❌ Error: DATABASE_URL is missing."
  echo "Usage: ./scripts/migrate_remote.sh <DATABASE_URL>"
  echo "Example: ./scripts/migrate_remote.sh 'postgres://adfluence:secret@dpg-xxxx.render.com/adfluence_db?sslmode=require'"
  exit 1
fi

if ! command -v psql >/dev/null 2>&1; then
  echo "❌ Error: psql command not found."
  echo "Please install PostgreSQL client tools or copy/paste 'migration/deploy_all.sql' directly into your cloud database SQL console."
  exit 1
fi

echo "🚀 Running all migrations against remote database..."
psql "$TARGET_URL" -v ON_ERROR_STOP=1 -f "$MIGRATION_FILE"

echo "✅ All migrations and seed data successfully applied!"
