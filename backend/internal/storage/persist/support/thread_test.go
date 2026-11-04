package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	"adinfluence_service_app/migration"
	"adinfluence_service_app/platform/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// threadDB returns a freshly migrated scratch database (needs MIGRATION_TEST_DSN).
func threadDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("MIGRATION_TEST_DSN not set")
	}
	root, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("support_test_%d", time.Now().UnixNano())
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

func newUser(t *testing.T, db *sql.DB, role string) string {
	t.Helper()
	var id string
	if err := db.QueryRow(`INSERT INTO users (email, password_hash, role, status) VALUES ($1, 'x', $2, 'active') RETURNING id`,
		fmt.Sprintf("%s%d@example.com", role, time.Now().UnixNano()), role).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func newTicket(t *testing.T, s *support, userID string) string {
	t.Helper()
	id := utils.GenerateUUID()
	if err := s.CreateTicket(context.Background(), businessDto.CreateSupportTicketInput{
		ID: id, UserID: userID, IssueType: "payment_question", Subject: "Payment not confirmed", Message: "I paid yesterday.", Priority: "normal",
	}); err != nil {
		t.Fatal(err)
	}
	return id
}

func TestSupportConversation(t *testing.T) {
	db := threadDB(t)
	s := &support{db: db}
	ctx := context.Background()
	user, admin, stranger := newUser(t, db, "business"), newUser(t, db, "admin"), newUser(t, db, "business")
	if _, err := db.Exec(`INSERT INTO business_profiles (user_id, owner_full_name, business_name) VALUES ($1, 'Hana', 'Addis Coffee Co.')`, user); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO admin_profiles (user_id, full_name) VALUES ($1, 'Meron')`, admin); err != nil {
		t.Fatal(err)
	}
	ticket := newTicket(t, s, user)
	status := func() string {
		var st string
		if err := db.QueryRow(`SELECT status::text FROM support_tickets WHERE id = $1`, ticket).Scan(&st); err != nil {
			t.Fatal(err)
		}
		return st
	}
	needsReply := func() int {
		stats, err := s.GetAdminTicketStats(ctx)
		if err != nil {
			t.Fatal(err)
		}
		return stats.NeedsReply
	}

	rec, err := s.GetTicketByUser(ctx, user, ticket)
	if err != nil || rec.LastReplyBy != "" {
		t.Fatalf("new ticket: %+v, %v", rec, err)
	}
	if needsReply() != 1 {
		t.Fatal("a new ticket needs a reply")
	}

	// An internal note alone is not an answer: the ticket stays open and the user doesn't see it.
	note := "Check the bank statement first."
	if err := s.SaveTicketReply(ctx, ticket, admin, nil, &note); err != nil {
		t.Fatal(err)
	}
	if status() != "open" {
		t.Fatalf("after internal note: status %s", status())
	}

	reply := "We confirmed your payment. The creator can start."
	if err := s.SaveTicketReply(ctx, ticket, admin, &reply, nil); err != nil {
		t.Fatal(err)
	}
	if status() != "in_review" {
		t.Fatalf("after reply: status %s", status())
	}
	msgs, err := s.ListMessages(ctx, ticket)
	if err != nil || len(msgs) != 1 || msgs[0].From != supportDto.FromSupport || msgs[0].SenderName != "Meron" || msgs[0].Body != reply {
		t.Fatalf("messages after reply: %+v, %v", msgs, err)
	}
	if needsReply() != 0 {
		t.Fatal("answered ticket still counted as needing a reply")
	}

	// The user replies after the ticket was resolved: it opens again.
	if err := s.UpdateTicketStatus(ctx, ticket, "resolved", admin); err != nil {
		t.Fatal(err)
	}
	msg, st, reopened, err := s.AddUserMessage(ctx, user, ticket, "  Thanks, but the creator still sees 'awaiting payment'.  ")
	if err != nil || !reopened || st != "open" || msg.From != supportDto.FromUser || msg.Body != "Thanks, but the creator still sees 'awaiting payment'." {
		t.Fatalf("user reply: %+v %s %v %v", msg, st, reopened, err)
	}
	var resolvedAt sql.NullTime
	db.QueryRow(`SELECT resolved_at FROM support_tickets WHERE id = $1`, ticket).Scan(&resolvedAt)
	if resolvedAt.Valid {
		t.Fatal("reopened ticket still has resolved_at")
	}
	if rec, _ := s.GetTicketByUser(ctx, user, ticket); rec.LastReplyBy != "user" {
		t.Fatalf("last reply by = %q, want user", rec.LastReplyBy)
	}
	if needsReply() != 1 {
		t.Fatal("user's reply should put the ticket back in the queue")
	}
	if msgs, _ := s.ListMessages(ctx, ticket); len(msgs) != 2 || msgs[1].SenderName != "Addis Coffee Co." {
		t.Fatalf("thread order/names: %+v", msgs)
	}

	// A second reply while it is open doesn't count as reopening.
	if _, st, reopened, err := s.AddUserMessage(ctx, user, ticket, "Any update?"); err != nil || reopened || st != "open" {
		t.Fatalf("second reply: %s %v %v", st, reopened, err)
	}

	// Someone else can't write on this ticket.
	if _, _, _, err := s.AddUserMessage(ctx, stranger, ticket, "hello"); !errors.Is(err, ErrTicketNotFound) {
		t.Fatalf("stranger: err = %v", err)
	}

	// The admin list shows who spoke last, newest activity first.
	older := newTicket(t, s, user)
	db.Exec(`UPDATE support_tickets SET created_at = NOW() - INTERVAL '2 days' WHERE id = $1`, older)
	rows, _, err := s.ListAdminTickets(ctx, adminDto.ListSupportTicketsQuery{})
	if err != nil || len(rows) != 2 || rows[0].ID != ticket || rows[0].LastReplyBy != "user" || rows[1].LastReplyBy != "" {
		t.Fatalf("admin list: %+v, %v", rows, err)
	}
	detail, err := s.GetAdminTicket(ctx, ticket)
	if err != nil || detail.LastReplyBy != "user" || detail.InternalNote == nil || *detail.InternalNote != note {
		t.Fatalf("admin detail: %+v, %v", detail, err)
	}
}

// Migration 00004 copies the old single admin reply into the conversation.
func TestLegacyReplyIsCopiedIntoThread(t *testing.T) {
	db := threadDB(t)
	s := &support{db: db}
	user, admin := newUser(t, db, "influencer"), newUser(t, db, "admin")
	ticket := newTicket(t, s, user)
	if _, err := db.Exec(`UPDATE support_tickets SET resolution_note = 'Your profile is approved.', assigned_admin_user_id = $2 WHERE id = $1`, ticket, admin); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile("../../../../migration/versions/00004_support_threads.sql")
	if err != nil {
		t.Fatal(err)
	}
	up := strings.Split(string(raw), "-- +goose Down")[0]
	for i := 0; i < 2; i++ { // running it again must not duplicate
		if _, err := db.Exec(up); err != nil {
			t.Fatal(err)
		}
	}
	msgs, err := s.ListMessages(context.Background(), ticket)
	if err != nil || len(msgs) != 1 || msgs[0].From != supportDto.FromSupport || msgs[0].Body != "Your profile is approved." {
		t.Fatalf("legacy reply: %+v, %v", msgs, err)
	}
}
