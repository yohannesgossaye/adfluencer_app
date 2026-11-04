package admin

import (
	"context"
	"errors"
	"testing"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
)

func TestAdminDisputeDetailAndSettledGuard(t *testing.T) {
	db := statusDB(t)
	a := &admin{db: db}
	ctx := context.Background()

	business, creator := insertUser(t, db, "business", "active"), insertUser(t, db, "influencer", "active")
	admin := insertUser(t, db, "admin", "active")
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`INSERT INTO business_profiles (user_id, owner_full_name, business_name) VALUES ($1, 'Owner', 'Shop')`, business)
	exec(`INSERT INTO influencer_profiles (user_id, full_name, creator_name) VALUES ($1, 'Selam', 'Selam')`, creator)
	var campaign, contract, dispute string
	if err := db.QueryRow(`INSERT INTO campaigns (business_user_id, title, description) VALUES ($1, 'Launch', 'Launch video') RETURNING id`, business).Scan(&campaign); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`INSERT INTO contracts (campaign_id, business_user_id, influencer_user_id, title, amount_etb, business_fee_amount_etb,
		influencer_fee_amount_etb, total_business_payable_etb, influencer_receivable_etb, platform_total_fee_etb, status)
		VALUES ($1, $2, $3, 'Launch video', 20000, 1000, 1000, 21000, 19000, 2000, 'disputed') RETURNING id`, campaign, business, creator).Scan(&contract); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`INSERT INTO disputes (contract_id, opened_by_user_id, against_user_id, issue_type, description, severity, status, reason)
		VALUES ($1, $2, $3, 'contract_issue', 'They asked for three extra videos.', 'medium', 'open', 'extra_work_requested') RETURNING id`,
		contract, creator, business).Scan(&dispute); err != nil {
		t.Fatal(err)
	}

	row, err := a.GetAdminDispute(ctx, dispute)
	if err != nil {
		t.Fatal(err)
	}
	if row.OpenedByRole != "influencer" || row.ContractStatus != "disputed" || row.Outcome != "" || row.IssueType != "extra_work_requested" {
		t.Fatalf("open dispute row = %+v", row)
	}
	if err := a.UpdateDisputeStatus(ctx, dispute, "under_review", admin); err != nil {
		t.Fatalf("mark under review: %v", err)
	}

	exec(`UPDATE disputes SET status = 'resolved', outcome = 'continue', resolved_at = NOW() WHERE id = $1`, dispute)
	if row, err = a.GetAdminDispute(ctx, dispute); err != nil || row.Outcome != "continue" {
		t.Fatalf("resolved row = %+v, err %v", row, err)
	}
	// A settled dispute can't be moved back to open: the contract has already moved on.
	if err := a.UpdateDisputeStatus(ctx, dispute, "open", admin); !errors.Is(err, contractDto.ErrDisputeAlreadyResolved) {
		t.Fatalf("reopen settled: err = %v", err)
	}
	if err := a.UpdateDisputeStatus(ctx, "00000000-0000-0000-0000-000000000000", "open", admin); !errors.Is(err, ErrAdminDisputeNotFound) {
		t.Fatalf("missing dispute: err = %v", err)
	}
}
