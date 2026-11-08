# Database backups

Every night at **02:00 Ethiopia time**, and before every deploy (Phase 4), the
database is dumped and uploaded to a private S3 bucket. S3 keeps backups for
**14 days**. The server keeps the last 3 days too, for a fast restore.

The site keeps running during a backup: `pg_dump` reads a consistent snapshot
while the database keeps serving users.

## Files

| File | Purpose |
|---|---|
| `backup.sh` | dump → check it's readable → upload to S3 → clean up old local copies |
| `restore.sh` | `verify`: restore into a throwaway database, compare with live, delete it · `replace-live`: emergency restore (keeps the old database, renamed) |
| `adfluence-backup.service` / `.timer` | systemd: run `backup.sh nightly` at 23:00 UTC (02:00 EAT) |
| `iam-policy.json` | what the server may do in S3: upload and read backups under `postgres/`, list them. **No delete.** |
| `backup.env.example` | settings → `/etc/adfluence/backup.env` (bucket name, region; no secrets) |

## How S3 access works (no keys stored anywhere)

The EC2 server gets an **IAM role** (a permission badge). AWS hands the server
short-lived credentials automatically through the instance metadata. The AWS
CLI runs as a container with `--network host` so it can reach them. Nothing
secret is written to disk.

The role can't delete backups, and the bucket has **versioning** on, so even
if the server were compromised, an attacker couldn't destroy the backups.

## Daily commands (on the server)

```bash
# Run a backup now
sudo /opt/adfluence/backup.sh manual

# When is the next nightly backup? Did the last one work?
systemctl list-timers adfluence-backup.timer
journalctl -u adfluence-backup.service --since "2 days ago" | grep -E "DONE|FAILED"

# Test that the newest backup restores completely (safe, live is untouched)
sudo /opt/adfluence/restore.sh verify latest

# Test a backup stored in S3
sudo /opt/adfluence/restore.sh verify postgres/2026/11/adfluence_db-20261108T230000Z-nightly.dump
```

## Emergency: restore the live database

Use this only when live data is damaged (a bad migration, mistaken deletes).
Everything written after the backup is lost.

1. Choose the backup: the newest **before** the damage.
   `ls -lt /var/backups/adfluence/`, or in S3: the AWS console → S3 → bucket → `postgres/`.
2. Run:
   ```bash
   sudo /opt/adfluence/restore.sh replace-live postgres/2026/11/<file>.dump
   ```
   It restores into a new database first (live is untouched while it runs),
   asks you to type `REPLACE LIVE`, then stops the API, swaps the databases,
   and starts the API.
3. The old database is **kept** as `adfluence_db_before_restore_<time>`. To undo,
   stop the API and rename the two back. When all is well, drop it:
   ```bash
   docker exec -it adfluence_postgres psql -U adfluence -d postgres -c 'DROP DATABASE "adfluence_db_before_restore_<time>"'
   ```

If Vault was also lost: restore Vault's volume first (unseal keys are in the
owner's password manager), or recreate the entry from the runbook in `../vault`.
