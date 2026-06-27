#!/bin/bash
# Runs once on first PostgreSQL container start (empty data volume).
set -euo pipefail

echo "==> Adfluence V1: applying schema migrations"

for file in auth admin campaigns marketplace support; do
  echo "    -> ${file}.sql"
  psql -v ON_ERROR_STOP=1 \
    --username "$POSTGRES_USER" \
    --dbname "$POSTGRES_DB" \
    -f "/docker-migrations/${file}.sql"
done

echo "==> Adfluence V1 schema complete (32 tables)"
