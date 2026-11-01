# Database migrations

Migrations run automatically when the backend starts (`MIGRATE_ON_START`,
default `true`). A deploy is: pull, rebuild, restart.

- `versions/` holds the migrations goose applies, in order, once each.
  Applied versions are recorded in the `goose_db_version` table.
- `versions/00001_baseline.sql` is the full schema (formerly `deploy_all.sql`).
  It runs only on an empty database. A database set up with the old
  hand-run scripts is recognised and marked as version 1 instead.
- The other `*.sql` files in this folder are the old hand-run migrations,
  kept for history. Do not add new ones here.

## Adding a change

1. Create `versions/000NN_short_name.sql` with the next number:

   ```sql
   -- +goose Up
   ALTER TABLE ... ;

   -- +goose Down
   ALTER TABLE ... ;
   ```

2. Never edit a file that has already been deployed; add a new one.
3. Wrap `DO $$ ... $$` blocks in `-- +goose StatementBegin` / `-- +goose StatementEnd`.

## Testing

`MIGRATION_TEST_DSN=postgres://user:pass@host:port/db go test ./migration/`
runs the migrations against throwaway databases on that server.
