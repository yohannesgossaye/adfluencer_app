package influencer

import (
	"testing"

	"adinfluence_service_app/internal/constant/model"
)

func TestCanUseSupport(t *testing.T) {
	allowed := []model.UserStatus{model.ActiveStatus, model.PendingProfileStatus, model.PendingApprovalStatus, model.RejectedStatus}
	refused := []model.UserStatus{model.BlockedStatus, model.SuspendedStatus, model.RemovedStatus}
	for _, s := range allowed {
		if !canUseSupport(s) {
			t.Errorf("%s should be able to contact support", s)
		}
	}
	for _, s := range refused {
		if canUseSupport(s) {
			t.Errorf("%s should not be able to contact support", s)
		}
	}
}
