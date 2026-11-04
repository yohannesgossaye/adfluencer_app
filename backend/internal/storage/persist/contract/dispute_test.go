package contract

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/migration"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// disputeDB returns a freshly migrated scratch database (needs MIGRATION_TEST_DSN).
func disputeDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("MIGRATION_TEST_DSN not set")
	}
	root, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("dispute_test_%d", time.Now().UnixNano())
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

type fixture struct{ business, creator, admin, contract string }

// activeContract creates a business, a creator, an admin and a paid, in-progress contract.
func activeContract(t *testing.T, db *sql.DB) fixture {
	t.Helper()
	var f fixture
	n := time.Now().UnixNano()
	user := func(role string) string {
		var id string
		if err := db.QueryRow(`INSERT INTO users (email, password_hash, role, status) VALUES ($1, 'x', $2, 'active') RETURNING id`,
			fmt.Sprintf("%s%d@example.com", role, n), role).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	f.business, f.creator, f.admin = user("business"), user("influencer"), user("admin")
	mustExec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	mustExec(`INSERT INTO business_profiles (user_id, owner_full_name, business_name) VALUES ($1, 'Owner', 'Shop')`, f.business)
	mustExec(`INSERT INTO influencer_profiles (user_id, full_name, creator_name) VALUES ($1, 'Selam', 'Selam')`, f.creator)
	var campaign string
	if err := db.QueryRow(`INSERT INTO campaigns (business_user_id, title, description, status) VALUES ($1, 'Launch', 'Launch video', 'active') RETURNING id`, f.business).Scan(&campaign); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`INSERT INTO contracts (campaign_id, business_user_id, influencer_user_id, title, amount_etb, business_fee_amount_etb,
		influencer_fee_amount_etb, total_business_payable_etb, influencer_receivable_etb, platform_total_fee_etb, status)
		VALUES ($1, $2, $3, 'Launch video', 20000, 1000, 1000, 21000, 19000, 2000, 'active') RETURNING id`, campaign, f.business, f.creator).Scan(&f.contract); err != nil {
		t.Fatal(err)
	}
	mustExec(`INSERT INTO payments (contract_id, business_user_id, influencer_user_id, amount_etb, business_fee_etb, influencer_fee_etb,
		platform_fee_total_etb, total_paid_by_business_etb, influencer_payout_etb, payment_method, status, paid_at)
		VALUES ($1, $2, $3, 20000, 1000, 1000, 2000, 21000, 19000, 'admin_confirmed', 'paid', NOW())`, f.contract, f.business, f.creator)
	return f
}

func contractStatus(t *testing.T, db *sql.DB, id string) string {
	var s string
	if err := db.QueryRow(`SELECT status::text FROM contracts WHERE id = $1`, id).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestDisputeFlow(t *testing.T) {
	db := disputeDB(t)
	repo := &contractRepo{db: db}
	ctx := context.Background()

	t.Run("open, block approval, one at a time, complete with payout", func(t *testing.T) {
		f := activeContract(t, db)
		d, err := repo.OpenDispute(ctx, f.business, f.contract, "not_delivered", "The video was never posted on TikTok.")
		if err != nil {
			t.Fatal(err)
		}
		if d.OpenedByRole != "business" || d.ReasonLabel != "Work not delivered" || d.Status != "open" {
			t.Fatalf("summary = %+v", d)
		}
		if s := contractStatus(t, db, f.contract); s != "disputed" {
			t.Fatalf("contract status = %s, want disputed", s)
		}
		// The payout can't be released while the dispute is open.
		if _, err := repo.ApproveWork(ctx, f.business, f.contract, "11111111-1111-1111-1111-111111111111"); !errors.Is(err, contractDto.ErrInvalidContractStatus) {
			t.Fatalf("approve while disputed: err = %v", err)
		}
		if _, err := repo.OpenDispute(ctx, f.creator, f.contract, "other", "Second report on the same contract."); !errors.Is(err, contractDto.ErrDisputeAlreadyOpen) {
			t.Fatalf("second dispute: err = %v", err)
		}
		cid, err := repo.ResolveDispute(ctx, f.admin, d.ID, contractDto.DisputeOutcomeComplete, "Creator posted it after all.")
		if err != nil || cid != f.contract {
			t.Fatalf("resolve: %v %s", err, cid)
		}
		if s := contractStatus(t, db, f.contract); s != "completed" {
			t.Fatalf("contract status = %s, want completed", s)
		}
		var payout float64
		if err := db.QueryRow(`SELECT p.amount_etb FROM payouts p JOIN payments pm ON pm.id = p.payment_id WHERE pm.contract_id = $1`, f.contract).Scan(&payout); err != nil || payout != 19000 {
			t.Fatalf("payout = %v, err %v", payout, err)
		}
		if _, err := repo.ResolveDispute(ctx, f.admin, d.ID, contractDto.DisputeOutcomeTerminate, "Trying to resolve twice."); !errors.Is(err, contractDto.ErrDisputeAlreadyResolved) {
			t.Fatalf("resolve twice: err = %v", err)
		}
	})

	t.Run("creator opens, admin continues, then terminates a new dispute", func(t *testing.T) {
		f := activeContract(t, db)
		d, err := repo.OpenDispute(ctx, f.creator, f.contract, "extra_work_requested", "They want three extra videos for free.")
		if err != nil || d.OpenedByRole != "influencer" {
			t.Fatalf("open: %v %+v", err, d)
		}
		if _, err := repo.ResolveDispute(ctx, f.admin, d.ID, contractDto.DisputeOutcomeContinue, "Business agreed to the original scope."); err != nil {
			t.Fatal(err)
		}
		if s := contractStatus(t, db, f.contract); s != "active" {
			t.Fatalf("after continue = %s, want active", s)
		}
		d2, err := repo.OpenDispute(ctx, f.business, f.contract, "creator_unresponsive", "No reply for two weeks after agreeing.")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := repo.ResolveDispute(ctx, f.admin, d2.ID, contractDto.DisputeOutcomeTerminate, "Creator unreachable; refunding the business."); err != nil {
			t.Fatal(err)
		}
		if s := contractStatus(t, db, f.contract); s != "cancelled" {
			t.Fatalf("after terminate = %s, want cancelled", s)
		}
	})

	t.Run("only parties, only in-progress contracts", func(t *testing.T) {
		f := activeContract(t, db)
		if _, err := repo.OpenDispute(ctx, f.admin, f.contract, "other", "I am not part of this contract."); !errors.Is(err, contractDto.ErrContractNotFound) {
			t.Fatalf("stranger: err = %v", err)
		}
		if _, err := db.Exec(`UPDATE contracts SET status = 'pending_influencer' WHERE id = $1`, f.contract); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.OpenDispute(ctx, f.business, f.contract, "other", "Contract not started yet."); !errors.Is(err, contractDto.ErrInvalidContractStatus) {
			t.Fatalf("not active: err = %v", err)
		}
	})
}

func TestOpenDisputeRequestValidate(t *testing.T) {
	ok := contractDto.OpenDisputeRequest{Reason: "not_delivered", Description: "The video was never posted."}
	if err := ok.Validate("business"); err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}
	if err := ok.Validate("influencer"); err == nil {
		t.Fatal("a creator can't pick a business-only reason")
	}
	if err := (contractDto.OpenDisputeRequest{Reason: "other", Description: "too short"}).Validate("influencer"); err == nil {
		t.Fatal("short description accepted")
	}
}
