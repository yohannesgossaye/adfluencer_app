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
	for _, c := range legacySchemaChecks {
		var one int
		if err := db.QueryRow(c.query).Scan(&one); err != nil {
			t.Errorf("fresh schema lacks %s", c.name)
		}
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

func TestLegacyDatabaseMissingMigrationIsRefused(t *testing.T) {
	db := scratchDB(t)
	if _, err := db.Exec(`CREATE TABLE users (id int)`); err != nil {
		t.Fatal(err)
	}
	_, err := Up(context.Background(), db)
	if err == nil || !strings.Contains(err.Error(), "missing earlier migrations") {
		t.Fatalf("err = %v, want a missing-migrations error", err)
	}
	var tracked bool
	db.QueryRow(`SELECT to_regclass('goose_db_version') IS NOT NULL`).Scan(&tracked)
	if tracked {
		t.Fatal("a refused database must not be marked as migrated")
	}
}
