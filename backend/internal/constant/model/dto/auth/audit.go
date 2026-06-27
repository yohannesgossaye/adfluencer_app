package auth

import "time"

type CreateAuditLogInput struct {
	ActorUserID  *string
	TargetUserID *string
	Action       string
	EntityType   string
	EntityID     *string
	Description  string
	IPAddress    string
	UserAgent    string
}

type LoginLockout struct {
	Email          string
	FailedAttempts int
	LockoutTier    int
	LockedUntil    *time.Time
}
