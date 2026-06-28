package module

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
)

type AuthService interface {
	SignUp(ctx context.Context, req authDto.SignUpRequest, meta authDto.SessionMeta) (*authDto.SignupResponse, error)
	Login(ctx context.Context, req authDto.LoginRequest, meta authDto.SessionMeta) (*authDto.LoginResponse, error)
	GetMe(ctx context.Context, userID string) (*authDto.MeResponse, error)
	Logout(ctx context.Context, userID string, req authDto.LogoutRequest, meta authDto.SessionMeta) error
	ForgotPassword(ctx context.Context, req authDto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) (*authDto.ResetPasswordResponse, error)
	SendVerifyEmail(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, userID string, req authDto.VerifyEmailRequest) (*authDto.VerifyEmailResponse, error)
	SendVerifyPhone(ctx context.Context, userID string) error
	VerifyPhone(ctx context.Context, userID string, req authDto.VerifyPhoneRequest) (*authDto.VerifyPhoneResponse, error)
	RefreshToken(ctx context.Context, req authDto.RefreshTokenRequest, meta authDto.SessionMeta) (*authDto.RefreshTokenResponse, error)
}

type InfluencerService interface {
	GetProfile(ctx context.Context, userID string) (*influencerDto.GetProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req influencerDto.UpdateProfileRequest, meta authDto.SessionMeta) (*influencerDto.GetProfileResponse, error)
	SubmitProfile(ctx context.Context, userID string, meta authDto.SessionMeta) (*influencerDto.SubmitProfileResponse, error)
	UploadAvatar(ctx context.Context, userID string, filename, contentType string, data []byte, meta authDto.SessionMeta) (*influencerDto.UploadAvatarResponse, error)
	UploadPortfolio(ctx context.Context, userID string, filename, contentType string, data []byte, meta authDto.SessionMeta) (*influencerDto.UploadPortfolioResponse, error)
	DeletePortfolio(ctx context.Context, userID, itemID string, meta authDto.SessionMeta) (*influencerDto.DeletePortfolioResponse, error)
	ListInvites(ctx context.Context, userID string, query influencerDto.ListInvitesQuery) (*influencerDto.ListInvitesResponse, error)
	GetInvite(ctx context.Context, userID, inviteID string) (*influencerDto.InviteDetailResponse, error)
	AcceptInvite(ctx context.Context, userID, inviteID string, meta authDto.SessionMeta) (*influencerDto.RespondInviteResponse, error)
	RejectInvite(ctx context.Context, userID, inviteID string, meta authDto.SessionMeta) (*influencerDto.RespondInviteResponse, error)
	ListContracts(ctx context.Context, userID string, query contractDto.ListQuery) (*influencerDto.ListContractsResponse, error)
	GetContract(ctx context.Context, userID, contractID string) (*contractDto.DetailResponse, error)
	AcceptContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*influencerDto.RespondContractResponse, error)
	RejectContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*influencerDto.RespondContractResponse, error)
	GetDashboardSummary(ctx context.Context, userID string) (*influencerDto.DashboardSummaryResponse, error)
	ListSupportInviteOptions(ctx context.Context, userID string) (*influencerDto.ListSupportInviteOptionsResponse, error)
	CreateSupportTicket(ctx context.Context, userID string, req influencerDto.CreateSupportTicketRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error)
	ListSupportTickets(ctx context.Context, userID string) (*influencerDto.ListSupportTicketsResponse, error)
	GetSupportTicket(ctx context.Context, userID, ticketID string) (*influencerDto.SupportTicketDetailResponse, error)
}

type AdminService interface {
	ListApplications(ctx context.Context, query adminDto.ListApplicationsQuery) (*adminDto.ListApplicationsResponse, error)
	GetApplication(ctx context.Context, userID string) (*adminDto.ApplicationDetailResponse, error)
	ApproveApplication(ctx context.Context, adminUserID, targetUserID string, req adminDto.ApproveApplicationRequest, meta authDto.SessionMeta) (*adminDto.ApproveApplicationResponse, error)
	RejectApplication(ctx context.Context, adminUserID, targetUserID string, req adminDto.RejectApplicationRequest, meta authDto.SessionMeta) (*adminDto.RejectApplicationResponse, error)
	UpdateAdminNotes(ctx context.Context, targetUserID string, req adminDto.UpdateAdminNotesRequest) error
	ListBusinessVerifications(ctx context.Context, query adminDto.ListBusinessVerificationsQuery) (*adminDto.ListBusinessVerificationsResponse, error)
	GetBusinessVerification(ctx context.Context, userID string) (*adminDto.BusinessVerificationDetailResponse, error)
	VerifyBusiness(ctx context.Context, adminUserID, targetUserID string, req adminDto.VerifyBusinessRequest, meta authDto.SessionMeta) (*adminDto.VerifyBusinessResponse, error)
	RevokeBusinessVerification(ctx context.Context, adminUserID, targetUserID string, req adminDto.RevokeBusinessVerificationRequest, meta authDto.SessionMeta) (*adminDto.RevokeBusinessVerificationResponse, error)
	ListSupportTickets(ctx context.Context, query adminDto.ListSupportTicketsQuery) (*adminDto.ListSupportTicketsResponse, error)
	GetSupportTicket(ctx context.Context, ticketID string) (*adminDto.SupportTicketDetailResponse, error)
	UpdateSupportTicketStatus(ctx context.Context, adminUserID, ticketID string, req adminDto.UpdateSupportTicketStatusRequest, meta authDto.SessionMeta) (*adminDto.UpdateSupportTicketStatusResponse, error)
	ReplySupportTicket(ctx context.Context, adminUserID, ticketID string, req adminDto.ReplySupportTicketRequest, meta authDto.SessionMeta) (*adminDto.ReplySupportTicketResponse, error)
	ListCategories(ctx context.Context) (*adminDto.ListCategoriesResponse, error)
	CreateCategory(ctx context.Context, adminUserID string, req adminDto.CreateCategoryRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error)
	UpdateCategory(ctx context.Context, adminUserID, categoryID string, req adminDto.UpdateCategoryRequest, meta authDto.SessionMeta) (*adminDto.UpdateCategoryResponse, error)
	DisableCategory(ctx context.Context, adminUserID, categoryID string, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error)
	EnableCategory(ctx context.Context, adminUserID, categoryID string, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error)
	ListUsers(ctx context.Context, query adminDto.ListUsersQuery) (*adminDto.ListUsersResponse, error)
	GetUser(ctx context.Context, userID string) (*adminDto.UserDetailResponse, error)
	BlockUser(ctx context.Context, adminUserID, targetUserID string, req adminDto.BlockUserRequest, meta authDto.SessionMeta) (*adminDto.BlockUserResponse, error)
	UnblockUser(ctx context.Context, adminUserID, targetUserID string, meta authDto.SessionMeta) (*adminDto.UnblockUserResponse, error)
	RemoveUser(ctx context.Context, adminUserID, targetUserID string, req adminDto.RemoveUserRequest, meta authDto.SessionMeta) (*adminDto.RemoveUserResponse, error)
	GetDashboardSummary(ctx context.Context) (*adminDto.DashboardSummaryResponse, error)
	ListAdminCampaigns(ctx context.Context, query adminDto.ListAdminCampaignsQuery) (*adminDto.ListAdminCampaignsResponse, error)
	GetAdminCampaign(ctx context.Context, campaignID string) (*adminDto.AdminCampaignDetailResponse, error)
	PauseAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error)
	CompleteAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error)
	FlagAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error)
	ListAdminDisputes(ctx context.Context, query adminDto.ListAdminDisputesQuery) (*adminDto.ListAdminDisputesResponse, error)
	GetAdminDispute(ctx context.Context, disputeID string) (*adminDto.AdminDisputeDetailResponse, error)
	UpdateDisputeStatus(ctx context.Context, adminUserID, disputeID string, req adminDto.UpdateDisputeStatusRequest, meta authDto.SessionMeta) (*adminDto.UpdateDisputeStatusResponse, error)
	UpdateDisputeNotes(ctx context.Context, adminUserID, disputeID string, req adminDto.UpdateDisputeNotesRequest, meta authDto.SessionMeta) (*adminDto.UpdateDisputeNotesResponse, error)
}

type BusinessService interface {
	GetProfile(ctx context.Context, userID string) (*businessDto.GetProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req businessDto.UpdateProfileRequest, meta authDto.SessionMeta) (*businessDto.GetProfileResponse, error)
	UploadLogo(ctx context.Context, userID string, filename, contentType string, data []byte, meta authDto.SessionMeta) (*businessDto.UploadLogoResponse, error)
	CompleteProfile(ctx context.Context, userID string, meta authDto.SessionMeta) (*businessDto.CompleteProfileResponse, error)
	CreateCampaign(ctx context.Context, userID string, req businessDto.CreateCampaignRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error)
	ListCampaigns(ctx context.Context, userID string, query businessDto.ListCampaignsQuery) (*businessDto.ListCampaignsResponse, error)
	GetCampaign(ctx context.Context, userID, campaignID string) (*businessDto.CampaignDetailResponse, error)
	UpdateCampaign(ctx context.Context, userID, campaignID string, req businessDto.UpdateCampaignRequest, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error)
	PublishCampaign(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error)
	CancelCampaign(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error)
	GenerateRecommendations(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.GenerateRecommendationsResponse, error)
	ListRecommendations(ctx context.Context, userID, campaignID string, query businessDto.ListRecommendationsQuery) (*businessDto.ListRecommendationsResponse, error)
	GetInfluencerForBusiness(ctx context.Context, userID, influencerUserID, campaignID string) (*businessDto.BusinessInfluencerViewResponse, error)
	CreateInvite(ctx context.Context, userID, campaignID string, req businessDto.CreateInviteRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error)
	ListInvites(ctx context.Context, userID, campaignID string) (*businessDto.ListInvitesResponse, error)
	CancelInvite(ctx context.Context, userID, campaignID, inviteID string, meta authDto.SessionMeta) error
	CreateContract(ctx context.Context, userID string, req businessDto.CreateContractRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error)
	ListContracts(ctx context.Context, userID string, query contractDto.ListQuery) (*businessDto.ListContractsResponse, error)
	GetContract(ctx context.Context, userID, contractID string) (*contractDto.DetailResponse, error)
	PayContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*businessDto.PayContractResponse, error)
	ApproveWork(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*businessDto.ApproveWorkResponse, error)
	CancelContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*businessDto.ContractMutationResponse, error)
	GetDashboardSummary(ctx context.Context, userID string) (*businessDto.DashboardSummaryResponse, error)
	ListSupportCampaignOptions(ctx context.Context, userID string) (*businessDto.ListSupportCampaignOptionsResponse, error)
	CreateSupportTicket(ctx context.Context, userID string, req businessDto.CreateSupportTicketRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error)
	ListSupportTickets(ctx context.Context, userID string) (*businessDto.ListSupportTicketsResponse, error)
	GetSupportTicket(ctx context.Context, userID, ticketID string) (*businessDto.SupportTicketDetailResponse, error)
}
