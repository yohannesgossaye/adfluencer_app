package admin

import (
	"context"
	"testing"
)

func TestDashboardNeedsActionCounts(t *testing.T) {
	db := statusDB(t)
	a := &admin{db: db}
	ctx := context.Background()

	biz := insertUser(t, db, "business", "active")
	if _, err := db.Exec(`INSERT INTO business_profiles (user_id, owner_full_name, business_name, verification_status)
		VALUES ($1, 'Owner', 'Shop', 'pending')`, biz); err != nil {
		t.Fatal(err)
	}
	for _, tk := range []struct{ priority, status string }{{"high", "open"}, {"urgent", "in_review"}, {"normal", "escalated"}, {"normal", "open"}, {"high", "resolved"}} {
		if _, err := db.Exec(`INSERT INTO support_tickets (created_by_user_id, subject, message, issue_type, priority, status)
			VALUES ($1, 's', 'm', 'general_feedback', $2::support_priority, $3::support_ticket_status)`, biz, tk.priority, tk.status); err != nil {
			t.Fatal(err)
		}
	}

	k, err := a.GetDashboardKPIs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if k.BusinessesToReview != 1 {
		t.Errorf("businesses to review = %d, want 1", k.BusinessesToReview)
	}
	if k.UrgentSupportTickets != 3 {
		t.Errorf("urgent tickets = %d, want 3 (high open, urgent in review, escalated)", k.UrgentSupportTickets)
	}
	if k.PaymentsToConfirm != 0 {
		t.Errorf("payments to confirm = %d, want 0", k.PaymentsToConfirm)
	}
}
