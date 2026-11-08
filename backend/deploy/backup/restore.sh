#!/usr/bin/env bash
# Restores an Adfluence backup.
#
#   restore.sh verify [SOURCE]         restore into a throwaway database, compare
#                                      every table with live, then delete it.
#                                      Safe: live data and the site are untouched.
#   restore.sh replace-live SOURCE     EMERGENCY: make SOURCE the live database.
#                                      The current database is kept, renamed, so
#                                      this can be undone.
#
# SOURCE: "latest" (newest backup on this server; the default), a local .dump
# file, or an S3 key such as postgres/2026/11/adfluence_db-20261108T230000Z-nightly.dump
set -Eeuo pipefail
umask 077

MODE=${1:-}
SOURCE=${2:-latest}
CONF=${BACKUP_CONFIG:-/etc/adfluence/backup.env}
# shellcheck source=/dev/null
[ -r "$CONF" ] && . "$CONF"

PG=${PG_CONTAINER:-adfluence_postgres}
API=${API_CONTAINER:-adfluence_backend}
DIR=${BACKUP_DIR:-/var/backups/adfluence}
REGION=${BACKUP_REGION:-eu-north-1}
AWS_IMAGE=${AWS_IMAGE:-amazon/aws-cli:2.22.35}

log() { printf '%s restore: %s\n' "$(date -u +%FT%TZ)" "$*"; }
# A temporary database left by a failed run is removed; live is never touched.
TEMP_DB=""
fail() {
  log "FAILED: $*"
  if [ -n "$TEMP_DB" ]; then
    docker exec "$PG" psql -U "$DB_USER" -d postgres -qc "DROP DATABASE IF EXISTS \"$TEMP_DB\"" </dev/null >/dev/null 2>&1 \
      && log "removed the temporary database $TEMP_DB"
  fi
  exit 1
}
trap 'fail "unexpected error on line $LINENO"' ERR

aws() {
  # Same user as the dump files (they are private, mode 600); HOME for the CLI's cache.
  # :z lets SELinux systems share the folder; elsewhere it does nothing.
  docker run --rm --network host --user "$(id -u):$(id -g)" -e HOME=/tmp -v "$DIR:$DIR:z" \
    ${AWS_ACCESS_KEY_ID:+-e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY} \
    ${AWS_SESSION_TOKEN:+-e AWS_SESSION_TOKEN} \
    "$AWS_IMAGE" --region "$REGION" ${S3_ENDPOINT:+--endpoint-url "$S3_ENDPOINT"} "$@"
}
# One SQL command. No -i: it must never read the terminal (that would swallow
# the answer to the REPLACE LIVE question).
psql_admin() { docker exec "$PG" psql -U "$DB_USER" -d postgres -v ON_ERROR_STOP=1 -qtA "$@" </dev/null; }

# Exact row count of every table, as "table count" lines, sorted.
table_counts() {
  docker exec -i "$PG" psql -U "$DB_USER" -d "$1" -v ON_ERROR_STOP=1 -qtA <<'SQL'
DO $$
DECLARE r record; n bigint;
BEGIN
  CREATE TEMP TABLE _counts (t text, n bigint);
  FOR r IN SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'public' LOOP
    EXECUTE format('SELECT count(*) FROM %I.%I', r.schemaname, r.tablename) INTO n;
    INSERT INTO _counts VALUES (r.tablename, n);
  END LOOP;
END $$;
SELECT t || ' ' || n FROM _counts ORDER BY t;
SQL
}

# Turns SOURCE into a local dump file path.
resolve_source() {
  case "$SOURCE" in
    latest)
      ls -1t "$DIR"/*.dump 2>/dev/null | head -1 ;;
    s3://* | postgres/*)
      [ -n "${BACKUP_BUCKET:-}" ] || fail "BACKUP_BUCKET is not set (in $CONF)"
      local key=${SOURCE#s3://"$BACKUP_BUCKET"/}
      mkdir -p "$DIR/downloads"
      local out
      out="$DIR/downloads/$(basename "$key")"
      aws s3 cp "s3://$BACKUP_BUCKET/$key" "$out" --only-show-errors >&2 || fail "download from S3 failed"
      echo "$out" ;;
    *)
      echo "$SOURCE" ;;
  esac
}

case "$MODE" in verify | replace-live) ;; *)
  echo "usage: restore.sh verify [SOURCE] | restore.sh replace-live SOURCE" >&2; exit 2 ;;
esac

DB_USER=${DB_USER:-$(docker exec "$PG" printenv POSTGRES_USER)}
DB_NAME=${DB_NAME:-$(docker exec "$PG" printenv POSTGRES_DB)}
FILE=$(resolve_source)
[ -n "$FILE" ] && [ -s "$FILE" ] || fail "no backup found for '$SOURCE'"
log "backup: $FILE"

STAMP=$(date -u +%Y%m%d%H%M%S)
TARGET="${DB_NAME}_restore_$STAMP"

log "restoring into a new database $TARGET (live $DB_NAME is not touched)"
psql_admin -c "CREATE DATABASE \"$TARGET\" OWNER \"$DB_USER\""
TEMP_DB=$TARGET
docker exec -i "$PG" pg_restore -U "$DB_USER" -d "$TARGET" --no-owner --exit-on-error <"$FILE" \
  || fail "pg_restore failed; nothing else changed"

LIVE=$(table_counts "$DB_NAME")
RESTORED=$(table_counts "$TARGET")
LIVE_TABLES=$(grep -c . <<<"$LIVE")
RESTORED_TABLES=$(grep -c . <<<"$RESTORED")
CHANGED=$(diff <(echo "$LIVE") <(echo "$RESTORED") | grep -c '^>' || true)
log "tables: live $LIVE_TABLES, backup $RESTORED_TABLES; tables whose row count changed since the backup: $CHANGED"
if [ "$CHANGED" -gt 0 ]; then
  join -a1 -a2 -e 0 -o 0,1.2,2.2 <(echo "$LIVE") <(echo "$RESTORED") | awk '$2 != $3 {print "  " $1 ": live " $2 ", backup " $3}'
fi

if [ "$MODE" = verify ]; then
  psql_admin -c "DROP DATABASE \"$TARGET\""
  TEMP_DB=""
  [ "$LIVE_TABLES" -eq "$RESTORED_TABLES" ] || fail "the backup has $RESTORED_TABLES tables, live has $LIVE_TABLES"
  log "RESTORE TEST PASSED: the backup restores completely ($RESTORED_TABLES tables); the test database was deleted"
  exit 0
fi

# replace-live
echo
echo "This makes the backup the LIVE database. Users lose everything written after it was taken."
answer=""
read -r -p "Type REPLACE LIVE to continue: " answer || true
[ "$answer" = "REPLACE LIVE" ] || fail "cancelled; nothing changed"
OLD="${DB_NAME}_before_restore_$STAMP"
log "stopping the API"
docker stop "$API" >/dev/null
docker exec -i "$PG" psql -U "$DB_USER" -d postgres -v ON_ERROR_STOP=1 -qtA >/dev/null <<SQL
SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname IN ('$DB_NAME', '$TARGET') AND pid <> pg_backend_pid();
ALTER DATABASE "$DB_NAME" RENAME TO "$OLD";
ALTER DATABASE "$TARGET" RENAME TO "$DB_NAME";
SQL
TEMP_DB=""
log "starting the API"
docker start "$API" >/dev/null
log "DONE: $DB_NAME now holds the backup. The previous database is kept as $OLD."
log "Undo: stop the API, rename $DB_NAME away and $OLD back, start the API."
log "When everything is fine: DROP DATABASE \"$OLD\" to free the space."
