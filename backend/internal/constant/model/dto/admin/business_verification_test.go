package admin

import "testing"

func TestRejectBusinessVerificationRequestValidate(t *testing.T) {
	if err := (RejectBusinessVerificationRequest{Reason: "TIN number does not match the trade license"}).Validate(); err != nil {
		t.Fatalf("valid reason rejected: %v", err)
	}
	for _, reason := range []string{"", "   ", "no"} {
		if err := (RejectBusinessVerificationRequest{Reason: reason}).Validate(); err == nil {
			t.Errorf("reason %q should be rejected: the business needs to know what to fix", reason)
		}
	}
}

func TestListBusinessVerificationsQueryNormalize(t *testing.T) {
	if got := (ListBusinessVerificationsQuery{}).Normalize().Status; got != "pending" {
		t.Fatalf("default status = %q, want the review queue", got)
	}
	for _, s := range []string{"pending", "verified", "rejected", "all"} {
		if got := (ListBusinessVerificationsQuery{Status: s}).Normalize().Status; got != s {
			t.Errorf("status %q normalized to %q", s, got)
		}
	}
}
