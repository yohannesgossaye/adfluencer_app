package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/module/influencer/helper"
	"adinfluence_service_app/internal/storage"
	platformStorage "adinfluence_service_app/platform/storage"
	"adinfluence_service_app/platform/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

type influencerModule struct {
	influencer storage.InfluencerPersistence
	contracts  storage.ContractPersistence
	support    storage.SupportPersistence
	auth       storage.AuthPersisitence
	files      platformStorage.Service
	logger     utils.Logger
}

func InitInfluencerModule(influencer storage.InfluencerPersistence, contracts storage.ContractPersistence, support storage.SupportPersistence, auth storage.AuthPersisitence, files platformStorage.Service, log utils.Logger) *influencerModule {
	if files == nil {
		files = platformStorage.NewNoop()
	}
	return &influencerModule{
		influencer: influencer,
		contracts:  contracts,
		support:    support,
		auth:       auth,
		files:      files,
		logger:     log,
	}
}

func (m *influencerModule) GetProfile(ctx context.Context, userID string, user *authDto.UserResponse) (*influencerDto.GetProfileResponse, error) {
	if err := ensureEditableStatus(user.Status); err != nil && user.Status != string(model.PendingApprovalStatus) && user.Status != string(model.ActiveStatus) {
		if user.Status == string(model.BlockedStatus) || user.Status == string(model.RemovedStatus) {
			return nil, err
		}
	}

	rec, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return m.buildResponse(user, rec), nil
}

func (m *influencerModule) UpdateProfile(ctx context.Context, userID string, req influencerDto.UpdateProfileRequest, meta authDto.SessionMeta) (*influencerDto.GetProfileResponse, error) {
	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}

	if err := ensureEditableStatus(user.Status); err != nil {
		return nil, err
	}

	categories, err := m.influencer.ResolveCategoryIDs(ctx, req.Niches)
	if err != nil {
		if errors.Is(err, influencerDto.ErrUnknownCategory) {
			return nil, customErr.ErrUnknownCategory
		}
		return nil, err
	}

	platforms := influencerDto.NormalizePlatforms(req.Platforms)
	socialMap := helper.SocialLinksToMap(req.SocialLinks)
	socialJSON, err := json.Marshal(socialMap)
	if err != nil {
		return nil, err
	}

	instagram, tiktok, youtube, total := helper.ResolveFollowers(helper.FollowerInput{
		Platforms:          platforms,
		Followers:          req.Followers,
		InstagramFollowers: req.InstagramFollowers,
		TikTokFollowers:    req.TikTokFollowers,
		YouTubeSubscribers: req.YouTubeSubscribers,
	})

	var region *string
	if v := strings.TrimSpace(req.Region); v != "" {
		region = &v
	}
	var audience *string
	if v := strings.TrimSpace(req.AudienceDescription); v != "" {
		audience = &v
	}

	categoryIDs := make([]string, 0, len(categories))
	for _, c := range categories {
		categoryIDs = append(categoryIDs, c.ID)
	}

	if err := m.influencer.UpdateProfile(ctx, influencerDto.UpdateProfileInput{
		UserID:              userID,
		CreatorName:         strings.TrimSpace(req.CreatorName),
		FullName:            strings.TrimSpace(req.CreatorName),
		Bio:                 strings.TrimSpace(req.Bio),
		City:                strings.TrimSpace(req.City),
		Region:              region,
		Languages:           req.Languages,
		SocialLinksJSON:     socialJSON,
		InstagramFollowers:  instagram,
		TikTokFollowers:     tiktok,
		YouTubeSubscribers:  youtube,
		TotalFollowers:      total,
		EngagementRate:      req.EngagementRate,
		MinPriceETB:         req.MinPriceETB,
		MaxPriceETB:         req.MaxPriceETB,
		AudienceDescription: audience,
		CategoryIDs:         categoryIDs,
		Platforms:           platforms,
	}); err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "influencer_profile",
		EntityID:     &userID,
		Description:  "influencer profile updated",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	rec, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	updatedUser, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return m.buildResponse(updatedUser, rec), nil
}

func (m *influencerModule) SubmitProfile(ctx context.Context, userID string, meta authDto.SessionMeta) (*influencerDto.SubmitProfileResponse, error) {
	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}

	switch model.UserStatus(user.Status) {
	case model.PendingProfileStatus, model.RejectedStatus:
	default:
		return nil, customErr.ErrInvalidStatusForSubmit
	}

	rec, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	completeness := helper.BuildCompleteness(helper.CompletenessInput{
		Profile:   rec.ProfileDetailResponse,
		Platforms: rec.Platforms,
	})
	if !completeness.CanSubmit {
		return nil, customErr.ErrProfileIncomplete
	}

	if err := m.influencer.SubmitProfile(ctx, userID); err != nil {
		if errors.Is(err, influencerDto.ErrInvalidStatusForSubmit) {
			return nil, customErr.ErrInvalidStatusForSubmit
		}
		return nil, err
	}

	now := time.Now().UTC()
	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "user",
		EntityID:     &userID,
		Description:  "influencer profile submitted for approval",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &influencerDto.SubmitProfileResponse{
		Status:      string(model.PendingApprovalStatus),
		Message:     "Profile submitted for review",
		SubmittedAt: now,
	}, nil
}

func (m *influencerModule) UploadAvatar(ctx context.Context, userID string, filename, contentType string, data []byte, meta authDto.SessionMeta) (*influencerDto.UploadAvatarResponse, error) {
	if !m.files.Enabled() {
		return nil, customErr.ErrStorageNotConfigured
	}

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}
	if err := ensureEditableStatus(user.Status); err != nil {
		return nil, err
	}

	validated, err := platformStorage.ValidateAvatarFile(filename, contentType, data)
	if err != nil {
		if errors.Is(err, platformStorage.ErrUnsupportedFormat) || errors.Is(err, platformStorage.ErrInvalidFile) {
			return nil, customErr.ErrInvalidFile
		}
		return nil, err
	}

	rec, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	objectKey := platformStorage.AvatarObjectKey(userID, validated.Extension)
	if _, err := m.files.UploadBytes(ctx, objectKey, bytes.NewReader(validated.Data), validated.Size, validated.ContentType); err != nil {
		return nil, err
	}

	publicURL, err := m.files.GetFileURL(ctx, objectKey, 0)
	if err != nil {
		_ = m.files.Delete(ctx, objectKey)
		return nil, err
	}

	if err := m.influencer.UpdateProfileImageURL(ctx, userID, publicURL); err != nil {
		_ = m.files.Delete(ctx, objectKey)
		return nil, err
	}

	if rec.ProfileImageURL != nil && strings.TrimSpace(*rec.ProfileImageURL) != "" {
		oldKey := platformStorage.ObjectKeyFromURL(*rec.ProfileImageURL, m.files.PublicURLBase())
		if oldKey != "" && oldKey != objectKey {
			if delErr := m.files.Delete(ctx, oldKey); delErr != nil {
				m.logger.Warn(ctx, "failed to delete old avatar", zap.String("object_key", oldKey), zap.String("reason", delErr.Error()))
			}
		}
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "influencer_profile",
		EntityID:     &userID,
		Description:  "influencer avatar uploaded",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &influencerDto.UploadAvatarResponse{ProfileImageURL: publicURL}, nil
}

func (m *influencerModule) UploadPortfolio(ctx context.Context, userID string, filename, contentType string, data []byte, meta authDto.SessionMeta) (*influencerDto.UploadPortfolioResponse, error) {
	if !m.files.Enabled() {
		return nil, customErr.ErrStorageNotConfigured
	}

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}
	if err := ensureEditableStatus(user.Status); err != nil {
		return nil, err
	}

	if _, err := m.influencer.GetProfileByUserID(ctx, userID); err != nil {
		return nil, err
	}

	validated, err := platformStorage.ValidatePortfolioFile(filename, contentType, data)
	if err != nil {
		if errors.Is(err, platformStorage.ErrUnsupportedFormat) || errors.Is(err, platformStorage.ErrInvalidFile) {
			return nil, customErr.ErrInvalidFile
		}
		return nil, err
	}

	count, err := m.influencer.CountPortfolioItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= 3 {
		return nil, customErr.ErrPortfolioLimitReached
	}

	objectKey := platformStorage.PortfolioObjectKey(userID, validated.Extension)
	if _, err := m.files.UploadBytes(ctx, objectKey, bytes.NewReader(validated.Data), validated.Size, validated.ContentType); err != nil {
		return nil, err
	}

	publicURL, err := m.files.GetFileURL(ctx, objectKey, 0)
	if err != nil {
		_ = m.files.Delete(ctx, objectKey)
		return nil, err
	}

	itemID := utils.GenerateUUID()
	item, err := m.influencer.CreatePortfolioItem(ctx, influencerDto.CreatePortfolioItemInput{
		ID:            itemID,
		UserID:        userID,
		FileURL:       publicURL,
		ObjectKey:     objectKey,
		MediaType:     validated.MediaType,
		MimeType:      validated.ContentType,
		FileSizeBytes: validated.Size,
	})
	if err != nil {
		_ = m.files.Delete(ctx, objectKey)
		if errors.Is(err, influencerDto.ErrPortfolioLimitReached) {
			return nil, customErr.ErrPortfolioLimitReached
		}
		return nil, err
	}

	items, err := m.influencer.ListPortfolioItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "influencer_portfolio_item",
		EntityID:     &itemID,
		Description:  "influencer portfolio sample uploaded",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &influencerDto.UploadPortfolioResponse{
		Item:           *item,
		PortfolioItems: items,
		PortfolioCount: len(items),
	}, nil
}

func (m *influencerModule) DeletePortfolio(ctx context.Context, userID, itemID string, meta authDto.SessionMeta) (*influencerDto.DeletePortfolioResponse, error) {
	if !m.files.Enabled() {
		return nil, customErr.ErrStorageNotConfigured
	}

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}
	if err := ensureEditableStatus(user.Status); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetPortfolioItem(ctx, userID, itemID)
	if err != nil {
		if errors.Is(err, influencerDto.ErrPortfolioItemNotFound) {
			return nil, customErr.ErrPortfolioItemNotFound
		}
		return nil, err
	}

	if err := m.influencer.DeletePortfolioItem(ctx, userID, itemID); err != nil {
		if errors.Is(err, influencerDto.ErrPortfolioItemNotFound) {
			return nil, customErr.ErrPortfolioItemNotFound
		}
		return nil, err
	}

	if delErr := m.files.Delete(ctx, rec.ObjectKey); delErr != nil {
		m.logger.Warn(ctx, "failed to delete portfolio object", zap.String("object_key", rec.ObjectKey), zap.String("reason", delErr.Error()))
	}

	items, err := m.influencer.ListPortfolioItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "influencer_portfolio_item",
		EntityID:     &itemID,
		Description:  "influencer portfolio sample deleted",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &influencerDto.DeletePortfolioResponse{
		PortfolioItems: items,
		PortfolioCount: len(items),
	}, nil
}

func (m *influencerModule) buildResponse(user *authDto.UserResponse, rec *influencerDto.ProfileRecord) *influencerDto.GetProfileResponse {
	detail := rec.ProfileDetailResponse
	if detail.PortfolioItems == nil {
		detail.PortfolioItems = []influencerDto.PortfolioItem{}
	}
	completeness := helper.BuildCompleteness(helper.CompletenessInput{
		Profile:   detail,
		Platforms: rec.Platforms,
	})

	return &influencerDto.GetProfileResponse{
		User: influencerDto.ProfileUserResponse{
			ID:            user.ID,
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			Status:        user.Status,
			StatusReason:  user.StatusReason,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
		},
		Profile:      detail,
		Completeness: completeness,
	}
}

func ensureEditableStatus(status string) error {
	switch model.UserStatus(status) {
	case model.PendingProfileStatus, model.PendingApprovalStatus, model.ActiveStatus, model.RejectedStatus:
		return nil
	case model.SuspendedStatus, model.BlockedStatus, model.RemovedStatus:
		return customErr.ErrProfileLocked
	default:
		return customErr.ErrProfileLocked
	}
}

func (m *influencerModule) writeAudit(ctx context.Context, input authDto.CreateAuditLogInput) {
	if err := m.influencer.CreateAuditLog(ctx, input); err != nil {
		m.logger.Warn(ctx, "audit log failed", zap.String("reason", err.Error()))
	}
}
