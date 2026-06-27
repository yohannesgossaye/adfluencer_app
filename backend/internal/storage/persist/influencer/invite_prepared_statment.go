package influencer

const (
	listInfluencerInvites = `
	SELECT
		ci.id,
		ci.status::text,
		ci.message,
		ci.created_at,
		ci.responded_at,
		c.id,
		c.title,
		c.status::text,
		c.budget_min_etb,
		c.budget_max_etb,
		c.start_date,
		c.end_date,
		bp.business_name
	FROM campaign_invites ci
	JOIN campaigns c ON c.id = ci.campaign_id
	JOIN business_profiles bp ON bp.user_id = ci.business_user_id
	WHERE ci.influencer_user_id = $1
		AND ($2 = '' OR ci.status::text = $2)
	ORDER BY ci.created_at DESC
	`

	getInfluencerInviteByID = `
	SELECT
		ci.id,
		ci.status::text,
		ci.message,
		ci.created_at,
		ci.responded_at,
		ci.business_user_id,
		c.id,
		c.title,
		c.status::text,
		c.budget_min_etb,
		c.budget_max_etb,
		c.start_date,
		c.end_date,
		bp.business_name
	FROM campaign_invites ci
	JOIN campaigns c ON c.id = ci.campaign_id
	JOIN business_profiles bp ON bp.user_id = ci.business_user_id
	WHERE ci.influencer_user_id = $1 AND ci.id = $2
	`

	getCampaignPlatformsForInvite = `
	SELECT platform::text
	FROM campaign_platforms
	WHERE campaign_id = $1
	ORDER BY platform
	`

	acceptInfluencerInvite = `
	UPDATE campaign_invites
	SET status = 'accepted', responded_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND influencer_user_id = $2 AND status = 'pending'
	`

	rejectInfluencerInvite = `
	UPDATE campaign_invites
	SET status = 'rejected', responded_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND influencer_user_id = $2 AND status = 'pending'
	`
)
