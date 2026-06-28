package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"errors"
)

func (m *influencerModule) ListInvites(ctx context.Context, userID string, query influencerDto.ListInvitesQuery) (*influencerDto.ListInvitesResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	items, err := m.influencer.ListInvites(ctx, userID, query)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []influencerDto.InviteListItem{}
	}

	return &influencerDto.ListInvitesResponse{Items: items}, nil
}

func (m *influencerModule) GetInvite(ctx context.Context, userID, inviteID string) (*influencerDto.InviteDetailResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}

	return inviteDetailFromRecord(rec), nil
}

func (m *influencerModule) AcceptInvite(ctx context.Context, userID, inviteID string, meta authDto.SessionMeta) (*influencerDto.RespondInviteResponse, error) {
	if err := m.requireVerifiedInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}
	if rec.Status != "pending" {
		return nil, customErr.ErrInvalidInviteStatus
	}
	if rec.Campaign.Status != "active" {
		return nil, customErr.ErrInvalidCampaignStatus
	}

	res, err := m.influencer.AcceptInvite(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.BusinessUserID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign_invite",
		EntityID:     &inviteID,
		Description:  "influencer accepted campaign invite",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return res, nil
}

func (m *influencerModule) RejectInvite(ctx context.Context, userID, inviteID string, meta authDto.SessionMeta) (*influencerDto.RespondInviteResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}
	if rec.Status != "pending" {
		return nil, customErr.ErrInvalidInviteStatus
	}

	res, err := m.influencer.RejectInvite(ctx, userID, inviteID)
	if err != nil {
		return nil, mapInviteError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.BusinessUserID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign_invite",
		EntityID:     &inviteID,
		Description:  "influencer rejected campaign invite",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

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

	rec, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, influencerDto.ErrProfileNotFound) {
			return customErr.ErrProfileNotFound
		}
		return err
	}
	if !rec.IsVerified {
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
