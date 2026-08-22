package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/module/business/helper"
	"adinfluence_service_app/internal/storage"
	platformStorage "adinfluence_service_app/platform/storage"
	"adinfluence_service_app/platform/utils"
	"bytes"
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

type businessModule struct {
	business   storage.BusinessPersistence
	influencer storage.InfluencerPersistence
	contracts  storage.ContractPersistence
	support    storage.SupportPersistence
	auth       storage.AuthPersisitence
	files      platformStorage.Service
	logger     utils.Logger
}

func InitBusinessModule(business storage.BusinessPersistence, influencer storage.InfluencerPersistence, contracts storage.ContractPersistence, support storage.SupportPersistence, auth storage.AuthPersisitence, files platformStorage.Service, log utils.Logger) *businessModule {
	if files == nil {
		files = platformStorage.NewNoop()
	}
	return &businessModule{
		business:   business,
		influencer: influencer,
		contracts:  contracts,
		support:    support,
		auth:       auth,
		files:      files,
		logger:     log,
	}
}

func (m *businessModule) GetProfile(ctx context.Context, userID string) (*businessDto.GetProfileResponse, error) {
	m.logger.Info(ctx, "fetching business profile", zap.String("user_id", userID))

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "business user not found", zap.String("user_id", userID))
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "failed to get business user from auth persistence", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}
	if err := ensureReadableStatus(user.Status); err != nil {
		m.logger.Warn(ctx, "business profile not readable due to status", zap.String("user_id", userID), zap.String("status", user.Status))
		return nil, err
	}

	rec, err := m.business.GetProfileByUserID(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to get business profile details", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	return m.buildResponse(user, rec), nil
}

func (m *businessModule) UpdateProfile(ctx context.Context, userID string, req businessDto.UpdateProfileRequest, meta authDto.SessionMeta) (*businessDto.GetProfileResponse, error) {
	m.logger.Info(ctx, "updating business profile", zap.String("user_id", userID), zap.String("business_name", req.BusinessName))

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "user not found during business profile update", zap.String("user_id", userID))
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "failed to query user for business update", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}
	if err := ensureEditableStatus(user.Status); err != nil {
		m.logger.Warn(ctx, "business profile is locked for edits", zap.String("user_id", userID), zap.String("status", user.Status))
		return nil, err
	}

	platforms := businessDto.NormalizePlatforms(req.PreferredPlatforms)

	rec, err := m.business.UpdateProfile(ctx, businessDto.UpdateProfileInput{
		UserID:                 userID,
		OwnerFullName:          strings.TrimSpace(req.OwnerFullName),
		BusinessName:           strings.TrimSpace(req.BusinessName),
		LegalBusinessName:      optionalString(req.LegalBusinessName),
		Description:            optionalString(req.Description),
		Industry:               optionalString(req.Industry),
		BusinessType:           optionalString(req.BusinessType),
		BusinessSize:           optionalString(req.BusinessSize),
		City:                   optionalString(req.City),
		Region:                 optionalString(req.Region),
		Address:                optionalString(req.Address),
		TINNumber:              optionalString(req.TINNumber),
		WebsiteURL:             optionalString(req.WebsiteURL),
		InstagramURL:           optionalString(req.SocialLinks.Instagram),
		TikTokURL:              optionalString(req.SocialLinks.TikTok),
		FacebookURL:            optionalString(req.SocialLinks.Facebook),
		LinkedInURL:            optionalString(req.SocialLinks.LinkedIn),
		PreferredPlatforms:     platforms,
		MonthlyMarketingBudget: req.MonthlyMarketingBudget,
	})
	if err != nil {
		m.logger.Error(ctx, "failed to update business profile in database", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "business_profile",
		EntityID:     &userID,
		Description:  "business profile updated",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "business profile updated successfully", zap.String("user_id", userID))

	return m.buildResponse(user, rec), nil
}

func (m *businessModule) UploadLogo(ctx context.Context, userID string, filename, contentType string, data []byte, meta authDto.SessionMeta) (*businessDto.UploadLogoResponse, error) {
	m.logger.Info(ctx, "uploading business logo", zap.String("user_id", userID), zap.String("filename", filename))

	if !m.files.Enabled() {
		m.logger.Warn(ctx, "storage service not configured for logo upload")
		return nil, customErr.ErrStorageNotConfigured
	}

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "failed to query user for logo upload", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}
	if err := ensureEditableStatus(user.Status); err != nil {
		m.logger.Warn(ctx, "logo upload rejected due to status", zap.String("user_id", userID), zap.String("status", user.Status))
		return nil, err
	}

	validated, err := platformStorage.ValidateAvatarFile(filename, contentType, data)
	if err != nil {
		if errors.Is(err, platformStorage.ErrUnsupportedFormat) || errors.Is(err, platformStorage.ErrInvalidFile) {
			m.logger.Warn(ctx, "logo validation failed: invalid format or file", zap.String("filename", filename))
			return nil, customErr.ErrInvalidFile
		}
		m.logger.Error(ctx, "logo validation error", zap.String("reason", err.Error()))
		return nil, err
	}

	rec, err := m.business.GetProfileByUserID(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to get business profile for logo upload", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	objectKey := platformStorage.BusinessLogoObjectKey(userID, validated.Extension)
	if _, err := m.files.UploadBytes(ctx, objectKey, bytes.NewReader(validated.Data), validated.Size, validated.ContentType); err != nil {
		m.logger.Error(ctx, "failed to upload logo bytes to S3", zap.String("object_key", objectKey), zap.String("reason", err.Error()))
		return nil, err
	}

	publicURL, err := m.files.GetFileURL(ctx, objectKey, 0)
	if err != nil {
		_ = m.files.Delete(ctx, objectKey)
		m.logger.Error(ctx, "failed to get public URL for logo", zap.String("object_key", objectKey), zap.String("reason", err.Error()))
		return nil, err
	}

	if err := m.business.UpdateLogoURL(ctx, userID, publicURL); err != nil {
		_ = m.files.Delete(ctx, objectKey)
		m.logger.Error(ctx, "failed to persist logo URL to database", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	if rec.LogoURL != nil && strings.TrimSpace(*rec.LogoURL) != "" {
		oldKey := platformStorage.ObjectKeyFromURL(*rec.LogoURL, m.files.PublicURLBase())
		if oldKey != "" && oldKey != objectKey {
			if delErr := m.files.Delete(ctx, oldKey); delErr != nil {
				m.logger.Warn(ctx, "failed to delete old logo", zap.String("object_key", oldKey), zap.String("reason", delErr.Error()))
			}
		}
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "business_profile",
		EntityID:     &userID,
		Description:  "business logo uploaded",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "business logo uploaded successfully", zap.String("user_id", userID), zap.String("url", publicURL))

	return &businessDto.UploadLogoResponse{LogoURL: publicURL}, nil
}

func (m *businessModule) CompleteProfile(ctx context.Context, userID string, meta authDto.SessionMeta) (*businessDto.CompleteProfileResponse, error) {
	m.logger.Info(ctx, "completing business profile", zap.String("user_id", userID))

	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "failed to query user for profile completion", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	switch model.UserStatus(user.Status) {
	case model.PendingProfileStatus, model.RejectedStatus:
	default:
		m.logger.Warn(ctx, "invalid status for profile complete", zap.String("user_id", userID), zap.String("status", user.Status))
		return nil, customErr.ErrInvalidStatusForComplete
	}

	rec, err := m.business.GetProfileByUserID(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to get business profile", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	completeness := helper.BuildCompleteness(rec.ProfileDetailResponse)
	if !completeness.CanComplete {
		m.logger.Warn(ctx, "business profile is incomplete", zap.String("user_id", userID))
		return nil, customErr.ErrProfileIncomplete
	}

	if err := m.business.CompleteProfile(ctx, userID); err != nil {
		if errors.Is(err, businessDto.ErrInvalidStatusForComplete) {
			return nil, customErr.ErrInvalidStatusForComplete
		}
		m.logger.Error(ctx, "failed to complete business profile in persistence", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	now := time.Now().UTC()
	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "user",
		EntityID:     &userID,
		Description:  "business profile completed — pending admin verification for marketplace",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "business profile completed successfully", zap.String("user_id", userID))

	return &businessDto.CompleteProfileResponse{
		Status:      string(model.ActiveStatus),
		Message:     "Business profile completed. Admin verification is required before marketplace visibility and contracts.",
		IsVerified:  rec.IsVerified,
		CompletedAt: now,
	}, nil
}

func (m *businessModule) buildResponse(user *authDto.UserResponse, rec *businessDto.ProfileRecord) *businessDto.GetProfileResponse {
	detail := rec.ProfileDetailResponse
	if detail.PreferredPlatforms == nil {
		detail.PreferredPlatforms = []string{}
	}
	completeness := helper.BuildCompleteness(detail)

	return &businessDto.GetProfileResponse{
		User: businessDto.ProfileUserResponse{
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
		MarketplaceAccess: helper.BuildMarketplaceAccess(detail, user.Status),
	}
}

func ensureEditableStatus(status string) error {
	switch model.UserStatus(status) {
	case model.PendingProfileStatus, model.ActiveStatus, model.RejectedStatus:
		return nil
	case model.SuspendedStatus, model.BlockedStatus, model.RemovedStatus:
		return customErr.ErrProfileLocked
	default:
		return customErr.ErrProfileLocked
	}
}

func ensureReadableStatus(status string) error {
	switch model.UserStatus(status) {
	case model.PendingProfileStatus, model.ActiveStatus, model.RejectedStatus, model.SuspendedStatus:
		return nil
	case model.BlockedStatus, model.RemovedStatus:
		return customErr.ErrProfileLocked
	default:
		return customErr.ErrProfileLocked
	}
}

func optionalString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

