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

// checkSchemaSQL lists every column, index, constraint and enum value from
// the baseline that the current database lacks (empty result = complete).
// Regenerate it whenever 00001_baseline.sql changes.
//
//go:embed check_schema.sql
var checkSchemaSQL string

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
// MissingSchema returns "kind name" for each baseline object the database lacks.
func MissingSchema(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, checkSchemaSQL)
	if err != nil {
		return nil, fmt.Errorf("check schema: %w", err)
	}
	defer rows.Close()
	var missing []string
	for rows.Next() {
		var problem, kind, name string
		if err := rows.Scan(&problem, &kind, &name); err != nil {
			return nil, fmt.Errorf("check schema: %w", err)
		}
		missing = append(missing, kind+" "+name)
	}
	return missing, rows.Err()
}

func adoptLegacyDatabase(ctx context.Context, db *sql.DB) error {
	var tracked, hasSchema bool
	if err := db.QueryRowContext(ctx, `SELECT to_regclass('goose_db_version') IS NOT NULL, to_regclass('users') IS NOT NULL`).Scan(&tracked, &hasSchema); err != nil {
		return fmt.Errorf("inspect database: %w", err)
	}
	if tracked || !hasSchema {
		return nil // already managed by goose, or empty and ready for the baseline
	}

	missing, err := MissingSchema(ctx, db)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		shown := missing
		if len(shown) > 20 {
			shown = append(shown[:20:20], fmt.Sprintf("... and %d more", len(missing)-20))
		}
		return fmt.Errorf("existing database is missing %d schema objects from earlier migrations; apply them before starting: %s", len(missing), strings.Join(shown, ", "))
	}

	if _, err := goose.EnsureDBVersionContext(ctx, db); err != nil {
		return fmt.Errorf("create goose_db_version: %w", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, true)`, baselineVersion); err != nil {
		return fmt.Errorf("mark baseline: %w", err)
	}
	return nil
}
