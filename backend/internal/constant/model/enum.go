package model

type Role string
type UserStatus string
type DocumentType string
type DocumentStatus string
type AuditAction string
type CampaignStatus string
type SocialPlatform string
type DeliverableType string
type CreatorSize string
type RecommendationStatus string
type InviteStatus string
type ConversationStatus string
type ConversationParticipantRole string
type MessageReportStatus string
type ContractStatus string
type ContractDeliverableStatus string
type PaymentStatus string
type PayoutStatus string
type SupportIssueType string
type SupportTicketStatus string
type SupportPriority string
type DisputeStatus string
type DisputeSeverity string
type NotificationType string

const (
	BusinessRole   Role = "business"
	InfluencerRole Role = "influencer"
	AdminRole      Role = "admin"
)

const (
	PendingProfileStatus  UserStatus = "pending_profile"
	PendingApprovalStatus UserStatus = "pending_approval"
	ActiveStatus          UserStatus = "active"
	RejectedStatus        UserStatus = "rejected"
	SuspendedStatus       UserStatus = "suspended"
	BlockedStatus         UserStatus = "blocked"
	RemovedStatus         UserStatus = "removed"
)

const (
	FaydaIDDocument         DocumentType = "fayda_id"
	PassportDocument        DocumentType = "passport"
	BusinessLicenseDocument DocumentType = "business_license"
	TradeLicenseDocument    DocumentType = "trade_license"
	TaxCertificateDocument  DocumentType = "tax_certificate"
	PortfolioSampleDocument DocumentType = "portfolio_sample"
	OtherDocument           DocumentType = "other"
)

const (
	PendingDocumentStatus  DocumentStatus = "pending"
	ApprovedDocumentStatus DocumentStatus = "approved"
	RejectedDocumentStatus DocumentStatus = "rejected"
)

const (
	LoginAction                AuditAction = "login"
	LogoutAction               AuditAction = "logout"
	SignupAction               AuditAction = "signup"
	PasswordResetAction        AuditAction = "password_reset"
	ProfileUpdateAction        AuditAction = "profile_update"
	StatusChangeAction         AuditAction = "status_change"
	ApproveAction              AuditAction = "approve"
	RejectAction               AuditAction = "reject"
	BlockAction                AuditAction = "block"
	UnblockAction              AuditAction = "unblock"
	SuspendAction              AuditAction = "suspend"
	UnsuspendAction            AuditAction = "unsuspend"
	DeleteAction               AuditAction = "delete"
	ContractStatusChangeAction AuditAction = "contract_status_change"
	PaymentStatusChangeAction  AuditAction = "payment_status_change"
	DisputeResolvedAction      AuditAction = "dispute_resolved"
	PermissionUpdateAction     AuditAction = "permission_update"
	AdminAction                AuditAction = "admin_action"
)

const (
	DraftCampaignStatus     CampaignStatus = "draft"
	PendingCampaignStatus   CampaignStatus = "pending"
	ActiveCampaignStatus    CampaignStatus = "active"
	CompletedCampaignStatus CampaignStatus = "completed"
	CancelledCampaignStatus CampaignStatus = "cancelled"
	PausedCampaignStatus    CampaignStatus = "paused"
)

const (
	TikTokPlatform    SocialPlatform = "tiktok"
	InstagramPlatform SocialPlatform = "instagram"
	YouTubePlatform   SocialPlatform = "youtube"
	TelegramPlatform  SocialPlatform = "telegram"
	FacebookPlatform  SocialPlatform = "facebook"
)

const (
	TikTokVideoDeliverable    DeliverableType = "tiktok_video"
	InstagramReelDeliverable  DeliverableType = "instagram_reel"
	InstagramStoryDeliverable DeliverableType = "instagram_story"
	YouTubeShortDeliverable   DeliverableType = "youtube_short"
	ProductReviewDeliverable  DeliverableType = "product_review"
	PostDeliverable           DeliverableType = "post"
	OtherDeliverable          DeliverableType = "other"
)

const (
	NanoCreator  CreatorSize = "nano"
	MicroCreator CreatorSize = "micro"
	MidCreator   CreatorSize = "mid"
	MacroCreator CreatorSize = "macro"
	MegaCreator  CreatorSize = "mega"
	AnyCreator   CreatorSize = "any"
)

const (
	RecommendedStatus RecommendationStatus = "recommended"
	InvitedStatus     RecommendationStatus = "invited"
	DismissedStatus   RecommendationStatus = "dismissed"
)

const (
	PendingInviteStatus   InviteStatus = "pending"
	AcceptedInviteStatus  InviteStatus = "accepted"
	RejectedInviteStatus  InviteStatus = "rejected"
	CancelledInviteStatus InviteStatus = "cancelled"
)

const (
	ActiveConversationStatus   ConversationStatus = "active"
	ArchivedConversationStatus ConversationStatus = "archived"
	ClosedConversationStatus   ConversationStatus = "closed"
)

const (
	BusinessParticipantRole   ConversationParticipantRole = "business"
	InfluencerParticipantRole ConversationParticipantRole = "influencer"
)

const (
	OpenMessageReportStatus        MessageReportStatus = "open"
	UnderReviewMessageReportStatus MessageReportStatus = "under_review"
	ResolvedMessageReportStatus    MessageReportStatus = "resolved"
	DismissedMessageReportStatus   MessageReportStatus = "dismissed"
)

const (
	DraftContractStatus             ContractStatus = "draft"
	PendingInfluencerContractStatus ContractStatus = "pending_influencer"
	PendingBusinessContractStatus   ContractStatus = "pending_business"
	AcceptedContractStatus          ContractStatus = "accepted"
	ActiveContractStatus            ContractStatus = "active"
	CompletedContractStatus         ContractStatus = "completed"
	CancelledContractStatus         ContractStatus = "cancelled"
	DisputedContractStatus          ContractStatus = "disputed"
	RejectedContractStatus          ContractStatus = "rejected"
)

const (
	PendingContractDeliverableStatus           ContractDeliverableStatus = "pending"
	SubmittedContractDeliverableStatus         ContractDeliverableStatus = "submitted"
	ApprovedContractDeliverableStatus          ContractDeliverableStatus = "approved"
	RevisionRequestedContractDeliverableStatus ContractDeliverableStatus = "revision_requested"
	CompletedContractDeliverableStatus         ContractDeliverableStatus = "completed"
)

const (
	PendingPaymentStatus   PaymentStatus = "pending"
	PaidPaymentStatus      PaymentStatus = "paid"
	FailedPaymentStatus    PaymentStatus = "failed"
	RefundedPaymentStatus  PaymentStatus = "refunded"
	CancelledPaymentStatus PaymentStatus = "cancelled"
)

const (
	PendingPayoutStatus    PayoutStatus = "pending"
	ProcessingPayoutStatus PayoutStatus = "processing"
	PaidPayoutStatus       PayoutStatus = "paid"
	FailedPayoutStatus     PayoutStatus = "failed"
	CancelledPayoutStatus  PayoutStatus = "cancelled"
)

const (
	CampaignIssue        SupportIssueType = "campaign_issue"
	ContractIssue        SupportIssueType = "contract_issue"
	PaymentQuestion      SupportIssueType = "payment_question"
	CreatorCommunication SupportIssueType = "creator_communication"
	ProfileApproval      SupportIssueType = "profile_approval"
	TechnicalIssue       SupportIssueType = "technical_issue"
	GeneralFeedback      SupportIssueType = "general_feedback"
	ComplaintIssue       SupportIssueType = "complaint"
)

const (
	OpenSupportTicketStatus      SupportTicketStatus = "open"
	InReviewSupportTicketStatus  SupportTicketStatus = "in_review"
	ResolvedSupportTicketStatus  SupportTicketStatus = "resolved"
	EscalatedSupportTicketStatus SupportTicketStatus = "escalated"
	ClosedSupportTicketStatus    SupportTicketStatus = "closed"
)

const (
	LowPriority    SupportPriority = "low"
	NormalPriority SupportPriority = "normal"
	HighPriority   SupportPriority = "high"
	UrgentPriority SupportPriority = "urgent"
)

const (
	OpenDisputeStatus        DisputeStatus = "open"
	UnderReviewDisputeStatus DisputeStatus = "under_review"
	ResolvedDisputeStatus    DisputeStatus = "resolved"
	EscalatedDisputeStatus   DisputeStatus = "escalated"
	ClosedDisputeStatus      DisputeStatus = "closed"
)

const (
	LowDisputeSeverity    DisputeSeverity = "low"
	MediumDisputeSeverity DisputeSeverity = "medium"
	HighDisputeSeverity   DisputeSeverity = "high"
)

const (
	CampaignInviteNotification       NotificationType = "campaign_invite"
	InviteAcceptedNotification       NotificationType = "invite_accepted"
	InviteRejectedNotification       NotificationType = "invite_rejected"
	ContractCreatedNotification      NotificationType = "contract_created"
	ContractAcceptedNotification     NotificationType = "contract_accepted"
	ContractCompletedNotification    NotificationType = "contract_completed"
	PaymentPaidNotification          NotificationType = "payment_paid"
	PayoutPaidNotification           NotificationType = "payout_paid"
	SupportReplyNotification         NotificationType = "support_reply"
	VerificationApprovedNotification NotificationType = "verification_approved"
	VerificationRejectedNotification NotificationType = "verification_rejected"
	ProfileApprovedNotification      NotificationType = "profile_approved"
	ProfileRejectedNotification      NotificationType = "profile_rejected"
	DisputeOpenedNotification        NotificationType = "dispute_opened"
	DisputeResolvedNotification      NotificationType = "dispute_resolved"
	MessageReportNotification        NotificationType = "message_report"
	GeneralNotification              NotificationType = "general"
)
