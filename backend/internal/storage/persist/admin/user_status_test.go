package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"adinfluence_service_app/migration"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// statusDB returns a freshly migrated scratch database. It needs a
// disposable Postgres server: set MIGRATION_TEST_DSN to a superuser DSN.
func statusDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("MIGRATION_TEST_DSN not set")
	}
	root, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("admin_status_test_%d", time.Now().UnixNano())
	if _, err := root.Exec("CREATE DATABASE " + name); err != nil {
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
		root.Exec("DROP DATABASE " + name + " WITH (FORCE)")
		root.Close()
	})
	if _, err := migration.Up(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return db
}

func insertUser(t *testing.T, db *sql.DB, role, status string) string {
	t.Helper()
	var id string
	err := db.QueryRow(`INSERT INTO users (email, password_hash, role, status)
		VALUES ($1, 'x', $2, $3) RETURNING id`,
		fmt.Sprintf("u%d@example.com", time.Now().UnixNano()), role, status).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestRestoreReturnsUserToPreviousStatus(t *testing.T) {
	db := statusDB(t)
	a := &admin{db: db}
	ctx := context.Background()

	cases := []struct {
		name   string
		start  string
		action func(id string) error
	}{
		{"removed", "pending_approval", func(id string) error { return a.RemoveUser(ctx, id, "spam") }},
		{"blocked", "active", func(id string) error { return a.BlockUser(ctx, id, "abuse") }},
		{"blocked then removed", "rejected", func(id string) error {
			if err := a.BlockUser(ctx, id, "abuse"); err != nil {
				return err
			}
			return a.RemoveUser(ctx, id, "spam")
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id := insertUser(t, db, "influencer", c.start)
			if err := c.action(id); err != nil {
				t.Fatal(err)
			}
			got, err := a.UnblockUser(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.start {
				t.Fatalf("restored to %q, want %q", got, c.start)
			}
		})
	}

	// Accounts locked before previous_status was recorded have none.
	t.Run("unknown previous status", func(t *testing.T) {
		for role, want := range map[string]string{"business": "active", "influencer": "pending_approval"} {
			id := insertUser(t, db, role, "suspended")
			got, err := a.UnblockUser(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("%s restored to %q, want %q", role, got, want)
			}
		}
	})

	t.Run("active user cannot be restored", func(t *testing.T) {
		id := insertUser(t, db, "business", "active")
		if _, err := a.UnblockUser(ctx, id); !errors.Is(err, ErrInvalidUserStatusChange) {
			t.Fatalf("err = %v, want ErrInvalidUserStatusChange", err)
		}
	})

	t.Run("admin cannot be blocked or removed", func(t *testing.T) {
		id := insertUser(t, db, "admin", "active")
		if err := a.BlockUser(ctx, id, "x"); !errors.Is(err, ErrCannotModifyAdminUser) {
			t.Fatalf("block err = %v", err)
		}
		if err := a.RemoveUser(ctx, id, "x"); !errors.Is(err, ErrCannotModifyAdminUser) {
			t.Fatalf("remove err = %v", err)
		}
	})
}
