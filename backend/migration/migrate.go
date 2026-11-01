// Package migration applies the versioned SQL files in versions/ with goose.
//
// The service runs Up on startup, so a deploy is: pull, rebuild, restart.
// Each file runs once, in order, and goose records it in goose_db_version.
//
// Adding a change: create versions/000NN_short_name.sql with
// "-- +goose Up" (and "-- +goose Down" when it can be undone). Never edit a
// file that has already been deployed; add a new one instead.
package migration

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
)

//go:embed versions/*.sql
var versions embed.FS

const baselineVersion = 1

// legacySchemaChecks are objects created by the last hand-run migrations
// (session_rotation, chat_integrity, business_logic, performance_indexes).
// A database created before goose must have all of them before it can be
// marked as being at the baseline.
var legacySchemaChecks = []struct{ name, query string }{
	{"user_sessions.refresh_ttl_seconds (session_rotation.sql)", `SELECT 1 FROM information_schema.columns WHERE table_name = 'user_sessions' AND column_name = 'refresh_ttl_seconds'`},
	{"messages.client_message_id (chat_integrity.sql)", `SELECT 1 FROM information_schema.columns WHERE table_name = 'messages' AND column_name = 'client_message_id'`},
	{"business_profiles.verification_status (business_logic.sql)", `SELECT 1 FROM information_schema.columns WHERE table_name = 'business_profiles' AND column_name = 'verification_status'`},
	{"users.previous_status (business_logic.sql)", `SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'previous_status'`},
	{"idx_contracts_business_status (performance_indexes.sql)", `SELECT 1 FROM pg_indexes WHERE indexname = 'idx_contracts_business_status'`},
}

// Up brings the database to the latest version and returns that version.
func Up(ctx context.Context, db *sql.DB) (int64, error) {
	goose.SetBaseFS(versions)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("postgres"); err != nil {
		return 0, err
	}

	if err := adoptLegacyDatabase(ctx, db); err != nil {
		return 0, err
	}
	if err := goose.UpContext(ctx, db, "versions"); err != nil {
		return 0, fmt.Errorf("apply migrations: %w", err)
	}
	return goose.GetDBVersionContext(ctx, db)
}

// adoptLegacyDatabase marks a database that was set up with the old
// hand-run scripts as being at the baseline, so the baseline (which would
// fail on existing tables) is not run against it.
func adoptLegacyDatabase(ctx context.Context, db *sql.DB) error {
	var tracked, hasSchema bool
	if err := db.QueryRowContext(ctx, `SELECT to_regclass('goose_db_version') IS NOT NULL, to_regclass('users') IS NOT NULL`).Scan(&tracked, &hasSchema); err != nil {
		return fmt.Errorf("inspect database: %w", err)
	}
	if tracked || !hasSchema {
		return nil // already managed by goose, or empty and ready for the baseline
	}

	var missing []string
	for _, c := range legacySchemaChecks {
		var one int
		err := db.QueryRowContext(ctx, c.query).Scan(&one)
		if err == sql.ErrNoRows {
			missing = append(missing, c.name)
			continue
		}
		if err != nil {
			return fmt.Errorf("inspect database: %w", err)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("existing database is missing earlier migrations, apply them before starting: %s", strings.Join(missing, ", "))
	}

	if _, err := goose.EnsureDBVersionContext(ctx, db); err != nil {
		return fmt.Errorf("create goose_db_version: %w", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, true)`, baselineVersion); err != nil {
		return fmt.Errorf("mark baseline: %w", err)
	}
	return nil
}
