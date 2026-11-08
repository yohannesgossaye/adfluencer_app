#!/usr/bin/env bash
# Backs up the Adfluence database to S3.
#
#   backup.sh [nightly|pre-deploy|manual]      (default: manual)
#
# 1. pg_dump inside the Postgres container (the site keeps running)
# 2. checks the file is a readable, complete dump
# 3. uploads it to s3://$BACKUP_BUCKET/postgres/YYYY/MM/<db>-<time>-<label>.dump
# 4. keeps the last few days on the server too, for a fast restore
#
# Settings: /etc/adfluence/backup.env (no secrets; S3 access comes from the
# server's IAM role). Exits non-zero on any failure, so systemd and deploy.sh notice.
set -Eeuo pipefail
umask 077

LABEL=${1:-manual}
CONF=${BACKUP_CONFIG:-/etc/adfluence/backup.env}
# shellcheck source=/dev/null
[ -r "$CONF" ] && . "$CONF"

PG=${PG_CONTAINER:-adfluence_postgres}
DIR=${BACKUP_DIR:-/var/backups/adfluence}
REGION=${BACKUP_REGION:-eu-north-1}
KEEP_LOCAL_DAYS=${KEEP_LOCAL_DAYS:-3}
AWS_IMAGE=${AWS_IMAGE:-amazon/aws-cli:2.22.35}

log() { printf '%s backup[%s]: %s\n' "$(date -u +%FT%TZ)" "$LABEL" "$*"; }
fail() { log "FAILED: $*"; rm -f "$DIR"/*.part; exit 1; }
trap 'fail "unexpected error on line $LINENO"' ERR

case "$LABEL" in nightly | pre-deploy | manual) ;; *) fail "label must be nightly, pre-deploy or manual" ;; esac
[ -n "${BACKUP_BUCKET:-}" ] || fail "BACKUP_BUCKET is not set (in $CONF)"

# The AWS CLI runs in a container. --network host lets it reach the EC2
# instance metadata, where AWS gives the IAM role's temporary credentials.
aws() {
  # Same user as the dump files (they are private, mode 600); HOME for the CLI's cache.
  # :z lets SELinux systems share the folder; elsewhere it does nothing.
  docker run --rm --network host --user "$(id -u):$(id -g)" -e HOME=/tmp -v "$DIR:$DIR:z" \
    ${AWS_ACCESS_KEY_ID:+-e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY} \
    ${AWS_SESSION_TOKEN:+-e AWS_SESSION_TOKEN} \
    "$AWS_IMAGE" --region "$REGION" ${S3_ENDPOINT:+--endpoint-url "$S3_ENDPOINT"} "$@"
}

mkdir -p "$DIR"
exec 9>"$DIR/.lock"
flock -n 9 || fail "another backup is already running"

[ "$(docker inspect -f '{{.State.Running}}' "$PG" 2>/dev/null)" = true ] || fail "the Postgres container $PG is not running"
DB_USER=${DB_USER:-$(docker exec "$PG" printenv POSTGRES_USER)}
DB_NAME=${DB_NAME:-$(docker exec "$PG" printenv POSTGRES_DB)}
[ -n "$DB_USER" ] && [ -n "$DB_NAME" ] || fail "could not find the database user/name in $PG"

STAMP=$(date -u +%Y%m%dT%H%M%SZ)
FILE="$DIR/$DB_NAME-$STAMP-$LABEL.dump"
log "dumping $DB_NAME from $PG"
docker exec "$PG" pg_dump -U "$DB_USER" -d "$DB_NAME" --format=custom --no-owner >"$FILE.part"
mv "$FILE.part" "$FILE"
[ -s "$FILE" ] || fail "the dump is empty"

# A dump that pg_restore can't read is not a backup.
TABLES=$(docker exec -i "$PG" pg_restore --list <"$FILE" | grep -c " TABLE DATA ") \
  || fail "the dump is not readable by pg_restore"
[ "$TABLES" -gt 0 ] || fail "the dump contains no table data"
SIZE=$(du -h "$FILE" | cut -f1)
log "dump OK: $TABLES tables, $SIZE"

KEY="postgres/$(date -u +%Y/%m)/$(basename "$FILE")"
log "uploading to s3://$BACKUP_BUCKET/$KEY"
# The bucket encrypts every object itself (default SSE-S3 encryption).
aws s3 cp "$FILE" "s3://$BACKUP_BUCKET/$KEY" --only-show-errors \
  || fail "upload to S3 failed (check the bucket name in $CONF and that the IAM role is attached to the server)"

# Keep only the last few days on the server; S3 keeps 14 days.
find "$DIR" -maxdepth 1 -name '*.dump' -mtime +"$KEEP_LOCAL_DAYS" -delete

# Optional: tell a monitor (Phase 5) that tonight's backup worked.
if [ "$LABEL" = nightly ] && [ -n "${HEALTHCHECK_URL:-}" ]; then
  curl -fsS -m 10 --retry 3 "$HEALTHCHECK_URL" >/dev/null || log "warning: could not ping the backup monitor"
fi

log "DONE: s3://$BACKUP_BUCKET/$KEY ($SIZE)"
