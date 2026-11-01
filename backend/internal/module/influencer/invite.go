package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"errors"

	"go.uber.org/zap"
)

func (m *influencerModule) ListInvites(ctx context.Context, userID string, query influencerDto.ListInvitesQuery) (*influencerDto.ListInvitesResponse, error) {
	m.logger.Info(ctx, "listing influencer invites", zap.String("user_id", userID), zap.String("status", query.Status))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	items, err := m.influencer.ListInvites(ctx, userID, query)
	if err != nil {
		m.logger.Error(ctx, "failed to query influencer invites", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}
	if items == nil {
		items = []influencerDto.InviteListItem{}
	}

	return &influencerDto.ListInvitesResponse{Items: items}, nil
}

func (m *influencerModule) GetInvite(ctx context.Context, userID, inviteID string) (*influencerDto.InviteDetailResponse, error) {
	m.logger.Info(ctx, "fetching influencer invite details", zap.String("user_id", userID), zap.String("invite_id", inviteID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		m.logger.Warn(ctx, "failed to get invite by ID", zap.String("invite_id", inviteID), zap.String("reason", err.Error()))
		return nil, mapInviteError(err)
	}

	return inviteDetailFromRecord(rec), nil
}

func (m *influencerModule) AcceptInvite(ctx context.Context, userID, inviteID string, meta authDto.SessionMeta) (*influencerDto.RespondInviteResponse, error) {
	m.logger.Info(ctx, "influencer accepting campaign invite", zap.String("user_id", userID), zap.String("invite_id", inviteID))

	if err := m.requireVerifiedInfluencer(ctx, userID); err != nil {
		m.logger.Warn(ctx, "cannot accept invite: influencer not verified", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	rec, err := m.influencer.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}
	if rec.Status != "pending" {
		m.logger.Warn(ctx, "cannot accept non-pending invite", zap.String("invite_id", inviteID), zap.String("status", rec.Status))
		return nil, customErr.ErrInvalidInviteStatus
	}
	if rec.Campaign.Status != "active" {
		m.logger.Warn(ctx, "cannot accept invite for non-active campaign", zap.String("campaign_id", rec.Campaign.ID), zap.String("campaign_status", rec.Campaign.Status))
		return nil, customErr.ErrInvalidCampaignStatus
	}

	res, err := m.influencer.AcceptInvite(ctx, userID, inviteID)
	if err != nil {
		m.logger.Error(ctx, "failed to accept invite in persistence", zap.String("invite_id", inviteID), zap.String("reason", err.Error()))
		return nil, mapInviteError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.BusinessUserID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign_invite",
		EntityID:     &inviteID,
		Description:  "influencer accepted campaign invite",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "influencer accepted invite successfully", zap.String("user_id", userID), zap.String("invite_id", inviteID))

	return res, nil
}

func (m *influencerModule) RejectInvite(ctx context.Context, userID, inviteID string, meta authDto.SessionMeta) (*influencerDto.RespondInviteResponse, error) {
	m.logger.Info(ctx, "influencer rejecting campaign invite", zap.String("user_id", userID), zap.String("invite_id", inviteID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}
	if rec.Status != "pending" {
		m.logger.Warn(ctx, "cannot reject non-pending invite", zap.String("invite_id", inviteID), zap.String("status", rec.Status))
		return nil, customErr.ErrInvalidInviteStatus
	}

	res, err := m.influencer.RejectInvite(ctx, userID, inviteID)
	if err != nil {
		m.logger.Error(ctx, "failed to reject invite in persistence", zap.String("invite_id", inviteID), zap.String("reason", err.Error()))
		return nil, mapInviteError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.BusinessUserID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign_invite",
		EntityID:     &inviteID,
		Description:  "influencer rejected campaign invite",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "influencer rejected invite successfully", zap.String("user_id", userID), zap.String("invite_id", inviteID))

	return res, nil
}

func inviteDetailFromRecord(rec *influencerDto.InviteRecord) *influencerDto.InviteDetailResponse {
	if rec.Campaign.Platforms == nil {
		rec.Campaign.Platforms = []string{}
	}
	return &influencerDto.InviteDetailResponse{
		ID:          rec.ID,
		Status:      rec.Status,
		Message:     rec.Message,
		Campaign:    rec.Campaign,
		CreatedAt:   rec.CreatedAt,
		RespondedAt: rec.RespondedAt,
	}
}

func (m *influencerModule) requireActiveInfluencer(ctx context.Context, userID string) error {
	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return customErr.ErrUnauthorized
		}
		return err
	}
	if model.UserStatus(user.Status) != model.ActiveStatus {
		return customErr.ErrAccountNotActive
	}
	return nil
}

func (m *influencerModule) requireVerifiedInfluencer(ctx context.Context, userID string) error {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return err
	}

	verified, err := m.influencer.IsProfileVerified(ctx, userID)
	if err != nil {
		if errors.Is(err, influencerDto.ErrProfileNotFound) {
			return customErr.ErrProfileNotFound
		}
		return err
	}
	if !verified {
		return customErr.ErrAccountNotVerified
	}
	return nil
}

func mapInviteError(err error) error {
	if errors.Is(err, influencerDto.ErrInviteNotFound) {
		return customErr.ErrInviteNotFound
	}
	if errors.Is(err, influencerDto.ErrInvalidInviteStatus) {
		return customErr.ErrInvalidInviteStatus
	}
	return err
}
