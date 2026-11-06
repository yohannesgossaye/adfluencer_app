package auth

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
