package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/module/admin/helper"
	influencerHelper "adinfluence_service_app/internal/module/influencer/helper"
	"adinfluence_service_app/internal/storage"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

type adminModule struct {
	admin      storage.AdminPersistence
	influencer storage.InfluencerPersistence
	support    storage.SupportPersistence
	logger     utils.Logger
}

func InitAdminModule(admin storage.AdminPersistence, influencer storage.InfluencerPersistence, support storage.SupportPersistence, log utils.Logger) *adminModule {
	return &adminModule{
		admin:      admin,
		influencer: influencer,
		support:    support,
		logger:     log,
	}
}

func (m *adminModule) ListApplications(ctx context.Context, query adminDto.ListApplicationsQuery) (*adminDto.ListApplicationsResponse, error) {
	query = query.Normalize()
	rows, total, err := m.admin.ListInfluencerApplications(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]adminDto.ApplicationListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapListItem(row))
	}

	return &adminDto.ListApplicationsResponse{
		Items: items,
		Pagination: adminDto.PaginationResponse{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	}, nil
}

func (m *adminModule) GetApplication(ctx context.Context, userID string) (*adminDto.ApplicationDetailResponse, error) {
	row, err := m.admin.GetInfluencerApplication(ctx, userID)
	if err != nil {
		if errors.Is(err, adminPersist.ErrApplicationNotFound) {
			return nil, customErr.ErrApplicationNotFound
		}
		return nil, err
	}

	profile, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	social := map[string]string{
		"tiktok":    profile.SocialLinks.TikTok,
		"instagram": profile.SocialLinks.Instagram,
		"youtube":   profile.SocialLinks.YouTube,
	}

	niches := make([]adminDto.NicheResponse, 0, len(profile.Niches))
	for _, n := range profile.Niches {
		niches = append(niches, adminDto.NicheResponse{ID: n.ID, Name: n.Name, Slug: n.Slug})
	}

	portfolio := make([]adminDto.PortfolioItem, 0, len(profile.PortfolioItems))
	for _, item := range profile.PortfolioItems {
		portfolio = append(portfolio, adminDto.PortfolioItem{
			ID:            item.ID,
			FileURL:       item.FileURL,
			MediaType:     item.MediaType,
			MimeType:      item.MimeType,
			FileSizeBytes: item.FileSizeBytes,
			SortOrder:     item.SortOrder,
			CreatedAt:     item.CreatedAt,
		})
	}

	completeness := influencerHelper.BuildCompleteness(influencerHelper.CompletenessInput{
		Profile:   profile.ProfileDetailResponse,
		Platforms: profile.Platforms,
	})

	adminNote := ""
	if row.ApprovalNote != nil {
		adminNote = *row.ApprovalNote
	}

	return &adminDto.ApplicationDetailResponse{
		UserID:       row.UserID,
		Status:       row.Status,
		StatusReason: row.StatusReason,
		Email:        row.Email,
		PhoneNumber:  row.PhoneNumber,
		Profile: adminDto.ApplicationProfileDetail{
			CreatorName:         profile.CreatorName,
			Bio:                 profile.Bio,
			ProfileImageURL:     profile.ProfileImageURL,
			City:                profile.City,
			Niches:              niches,
			Platforms:           profile.Platforms,
			SocialLinks:         social,
			TotalFollowers:      profile.TotalFollowers,
			EngagementRate:      profile.EngagementRate,
			MinPriceETB:         profile.MinPriceETB,
			MaxPriceETB:         profile.MaxPriceETB,
			AudienceDescription: profile.AudienceDescription,
			Languages:           profile.Languages,
			PortfolioItems:      portfolio,
		},
		CompletenessPercent: completeness.Percent,
		AdminNote:           adminNote,
		SubmittedAt:         row.UpdatedAt,
	}, nil
}

func (m *adminModule) ApproveApplication(ctx context.Context, adminUserID, targetUserID string, req adminDto.ApproveApplicationRequest, meta authDto.SessionMeta) (*adminDto.ApproveApplicationResponse, error) {
	var note *string
	if v := strings.TrimSpace(req.AdminNote); v != "" {
		note = &v
	}

	if err := m.admin.ApproveInfluencerApplication(ctx, targetUserID, adminUserID, note); err != nil {
		if errors.Is(err, adminPersist.ErrInvalidStatusForReview) {
			return nil, customErr.ErrInvalidStatusForReview
		}
		return nil, err
	}

	now := time.Now().UTC()
	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.ApproveAction),
		EntityType:   "user",
		EntityID:     &targetUserID,
		Description:  "influencer application approved",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &adminDto.ApproveApplicationResponse{
		UserID:     targetUserID,
		Status:     string(model.ActiveStatus),
		ApprovedAt: now,
	}, nil
}

func (m *adminModule) RejectApplication(ctx context.Context, adminUserID, targetUserID string, req adminDto.RejectApplicationRequest, meta authDto.SessionMeta) (*adminDto.RejectApplicationResponse, error) {
	reason := strings.TrimSpace(req.Reason)
	var note *string
	if v := strings.TrimSpace(req.AdminNote); v != "" {
		note = &v
	}

	if err := m.admin.RejectInfluencerApplication(ctx, targetUserID, reason, note); err != nil {
		if errors.Is(err, adminPersist.ErrInvalidStatusForReview) {
			return nil, customErr.ErrInvalidStatusForReview
		}
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.RejectAction),
		EntityType:   "user",
		EntityID:     &targetUserID,
		Description:  "influencer application rejected",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &adminDto.RejectApplicationResponse{
		UserID: targetUserID,
		Status: string(model.RejectedStatus),
		Reason: reason,
	}, nil
}

func (m *adminModule) UpdateAdminNotes(ctx context.Context, targetUserID string, req adminDto.UpdateAdminNotesRequest) error {
	if err := m.admin.UpdateInfluencerAdminNote(ctx, targetUserID, strings.TrimSpace(req.AdminNote)); err != nil {
		if errors.Is(err, adminPersist.ErrApplicationNotFound) {
			return customErr.ErrApplicationNotFound
		}
		return err
	}
	return nil
}

func mapListItem(row adminDto.ApplicationRow) adminDto.ApplicationListItem {
	return adminDto.ApplicationListItem{
		UserID:              row.UserID,
		CreatorName:         row.CreatorName,
		ProfileImageURL:     row.ProfileImageURL,
		PrimaryNiche:        row.PrimaryNiche,
		PrimaryPlatform:     row.PrimaryPlatform,
		TotalFollowers:      row.TotalFollowers,
		EngagementRate:      row.EngagementRate,
		CompletenessPercent: helper.CompletenessFromRow(row),
		PricingSummary:      helper.PricingSummary(row.MinPriceETB, row.MaxPriceETB),
		Status:              row.Status,
		SubmittedAt:         row.UpdatedAt,
		BioPreview:          helper.BioPreview(row.Bio),
	}
}

func (m *adminModule) writeAudit(ctx context.Context, input authDto.CreateAuditLogInput) {
	if err := m.admin.CreateAuditLog(ctx, input); err != nil {
		m.logger.Warn(ctx, "audit log failed", zap.String("reason", err.Error()))
	}
}
