package business

const (
	insertCampaignInvite = `
	INSERT INTO campaign_invites (
		id, campaign_id, business_user_id, influencer_user_id, message, status
	) VALUES ($1, $2, $3, $4, $5, 'pending')
	`

	listCampaignInvites = `
	SELECT
		ci.id,
		ci.campaign_id,
		ci.influencer_user_id,
		ip.creator_name,
		ci.status::text,
		ci.message,
		ci.created_at
	FROM campaign_invites ci
	JOIN influencer_profiles ip ON ip.user_id = ci.influencer_user_id
	WHERE ci.campaign_id = $1 AND ci.business_user_id = $2
	ORDER BY ci.created_at DESC
	`

	getCampaignInviteByID = `
	SELECT
		ci.id,
		ci.campaign_id,
		ci.business_user_id,
		ci.influencer_user_id,
		ip.creator_name,
		ci.status::text,
		ci.message,
		ci.created_at
	FROM campaign_invites ci
	JOIN influencer_profiles ip ON ip.user_id = ci.influencer_user_id
	WHERE ci.id = $1 AND ci.campaign_id = $2 AND ci.business_user_id = $3
	`

	cancelPendingInvite = `
	UPDATE campaign_invites
	SET status = 'cancelled', updated_at = NOW()
	WHERE id = $1 AND campaign_id = $2 AND business_user_id = $3 AND status = 'pending'
	`

	updateRecommendationStatusInvited = `
	UPDATE campaign_recommendations
	SET status = 'invited'
	WHERE campaign_id = $1 AND influencer_user_id = $2
	`

	inviteExistsForCampaign = `
	SELECT EXISTS (
		SELECT 1
		FROM campaign_invites
		WHERE campaign_id = $1
			AND influencer_user_id = $2
			AND status IN ('pending', 'accepted')
	)
	`
)
