package helper

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"testing"
)

func TestBuildMarketplaceAccess(t *testing.T) {
	profile := func(verified bool, status string) businessDto.ProfileDetailResponse {
		return businessDto.ProfileDetailResponse{
			IsVerified:   verified,
			Verification: businessDto.VerificationInfo{Status: status},
		}
	}

	cases := []struct {
		name       string
		profile    businessDto.ProfileDetailResponse
		userStatus string
		canAct     bool
	}{
		{"new business", profile(false, businessDto.VerificationUnsubmitted), "pending_profile", false},
		{"awaiting review", profile(false, businessDto.VerificationPending), "active", false},
		{"rejected", profile(false, businessDto.VerificationRejected), "active", false},
		{"verified", profile(true, businessDto.VerificationVerified), "active", true},
		{"verified but blocked", profile(true, businessDto.VerificationVerified), "blocked", false},
	}
	for _, tc := range cases {
		got := BuildMarketplaceAccess(tc.profile, tc.userStatus)
		if got.CanPublishCampaigns != tc.canAct || got.CanCreateContracts != tc.canAct || got.VisibleInMarketplace != tc.canAct {
			t.Errorf("%s: access = %+v, want canAct=%v", tc.name, got, tc.canAct)
		}
		if got.Message == "" {
			t.Errorf("%s: missing message for the business", tc.name)
		}
		if got.VerificationStatus != tc.profile.Verification.Status {
			t.Errorf("%s: verification status = %q", tc.name, got.VerificationStatus)
		}
	}
}
