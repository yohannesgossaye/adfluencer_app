package admin

import "testing"

func TestListApplicationsQueryMapsTabsToUserStatus(t *testing.T) {
	cases := map[string]string{"pending": "pending_approval", "approved": "active", "rejected": "rejected", "all": "", "": "", " Pending ": "pending_approval"}
	for in, want := range cases {
		if got := (ListApplicationsQuery{Status: in}).Normalize().Status; got != want {
			t.Errorf("Normalize(%q).Status = %q, want %q", in, got, want)
		}
	}
}
