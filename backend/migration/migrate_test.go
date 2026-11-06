package migration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// These tests need a disposable Postgres server. Set MIGRATION_TEST_DSN to
// a superuser DSN; each test creates and drops its own database.
func scratchDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("MIGRATION_TEST_DSN not set")
	}
	admin, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("goose_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + name); err != nil {
		t.Fatal(err)
	}
	i := strings.LastIndex(dsn, "/")
	q := ""
	if j := strings.Index(dsn[i:], "?"); j >= 0 {
		q = dsn[i+j:]
	}
	db, err := sql.Open("pgx", dsn[:i+1]+name+q)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
		admin.Exec("DROP DATABASE " + name + " WITH (FORCE)")
		admin.Close()
	})
	return db
}

func latest(t *testing.T) int64 {
	entries, err := versions.ReadDir("versions")
	if err != nil {
		t.Fatal(err)
	}
	return int64(len(entries))
}

func TestFreshDatabaseGetsFullSchema(t *testing.T) {
	db := scratchDB(t)
	ctx := context.Background()
	v, err := Up(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if v != latest(t) {
		t.Fatalf("version = %d, want %d", v, latest(t))
	}
	missing, err := MissingSchema(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	// check_schema.sql describes the baseline; later migrations may replace
	// a baseline object on purpose.
	replacedLater := map[string]bool{
		"index uq_conversations_pair_campaign": true, // 00006: one chat per pair
	}
	var unexpected []string
	for _, m := range missing {
		if !replacedLater[m] {
			unexpected = append(unexpected, m)
		}
	}
	if len(unexpected) > 0 {
		t.Fatalf("fresh schema: missing %v", unexpected)
	}
	// Running again is a no-op.
	if v2, err := Up(ctx, db); err != nil || v2 != v {
		t.Fatalf("second Up: version %d, err %v", v2, err)
	}
}

func TestLegacyDatabaseIsAdoptedWithoutRerunningBaseline(t *testing.T) {
	db := scratchDB(t)
	ctx := context.Background()
	// Build a "hand-migrated" database: the baseline SQL without goose.
	raw, err := versions.ReadFile("versions/00001_baseline.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(raw)); err != nil {
		t.Fatal(err)
	}
	v, err := Up(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if v != latest(t) {
		t.Fatalf("version = %d, want %d", v, latest(t))
	}
}

func TestLegacyDatabaseWithOneMissingColumnIsRefused(t *testing.T) {
	db := scratchDB(t)
	raw, err := versions.ReadFile("versions/00001_baseline.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(raw)); err != nil {
		t.Fatal(err)
	}
	// The exact gap production had: chat_integrity.sql never applied.
	if _, err := db.Exec(`ALTER TABLE messages DROP COLUMN client_message_id CASCADE`); err != nil {
		t.Fatal(err)
	}
	_, err = Up(context.Background(), db)
	if err == nil || !strings.Contains(err.Error(), "column messages.client_message_id") || !strings.Contains(err.Error(), "index uq_messages_sender_client_message_id") {
		t.Fatalf("err = %v, want both missing objects named", err)
	}
}

func TestLegacyDatabaseMissingMigrationIsRefused(t *testing.T) {
	db := scratchDB(t)
	if _, err := db.Exec(`CREATE TABLE users (id int)`); err != nil {
		t.Fatal(err)
	}
	_, err := Up(context.Background(), db)
	if err == nil || !strings.Contains(err.Error(), "is missing") || !strings.Contains(err.Error(), "more") {
		t.Fatalf("err = %v, want a missing-migrations error", err)
	}
	var tracked bool
	db.QueryRow(`SELECT to_regclass('goose_db_version') IS NOT NULL`).Scan(&tracked)
	if tracked {
		t.Fatal("a refused database must not be marked as migrated")
	}
}
