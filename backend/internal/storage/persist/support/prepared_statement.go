package support

const (
	listBusinessCampaignOptions = `
	SELECT id, title
	FROM campaigns
	WHERE business_user_id = $1
	ORDER BY updated_at DESC
	LIMIT 50
	`

	insertSupportTicket = `
	INSERT INTO support_tickets (
		id, created_by_user_id, subject, message, issue_type, priority, status, related_campaign_id
	) VALUES ($1, $2, $3, $4, $5::support_issue_type, $6::support_priority, 'open', $7)
	`

	listSupportTicketsByUser = `
	SELECT
		st.id,
		st.subject,
		st.message,
		st.issue_type::text,
		st.priority::text,
		st.status::text,
		st.resolution_note,
		st.related_campaign_id,
		c.title,
		st.created_at
	FROM support_tickets st
	LEFT JOIN campaigns c ON c.id = st.related_campaign_id
	WHERE st.created_by_user_id = $1
	ORDER BY st.created_at DESC
	`

	getSupportTicketByUser = `
	SELECT
		st.id,
		st.subject,
		st.message,
		st.issue_type::text,
		st.priority::text,
		st.status::text,
		st.resolution_note,
		st.related_campaign_id,
		c.title,
		st.created_at
	FROM support_tickets st
	LEFT JOIN campaigns c ON c.id = st.related_campaign_id
	WHERE st.id = $1 AND st.created_by_user_id = $2
	`

	businessOwnsCampaign = `
	SELECT EXISTS (
		SELECT 1 FROM campaigns WHERE id = $1 AND business_user_id = $2
	)
	`

	listAdminSupportTickets = `
	SELECT
		st.id,
		st.subject,
		st.message,
		st.issue_type::text,
		st.priority::text,
		st.status::text,
		st.resolution_note,
		u.id,
		COALESCE(bp.business_name, ip.creator_name, split_part(u.email, '@', 1)),
		u.email,
		u.role::text,
		st.related_campaign_id,
		c.title,
		st.created_at,
		st.updated_at,
		COUNT(*) OVER()::int AS total_count
	FROM support_tickets st
	INNER JOIN users u ON u.id = st.created_by_user_id
	LEFT JOIN business_profiles bp ON bp.user_id = u.id AND u.role = 'business'
	LEFT JOIN influencer_profiles ip ON ip.user_id = u.id AND u.role = 'influencer'
	LEFT JOIN campaigns c ON c.id = st.related_campaign_id
	WHERE ($1 = '' OR u.role::text = $1)
		AND ($2 = '' OR st.status::text = $2)
		AND ($3 = '' OR st.priority::text = $3)
		AND (
			$4 = ''
			OR st.subject ILIKE '%' || $4 || '%'
			OR st.message ILIKE '%' || $4 || '%'
			OR u.email ILIKE '%' || $4 || '%'
		)
	ORDER BY st.created_at DESC
	LIMIT $5 OFFSET $6
	`

	getAdminSupportTicketStats = `
	SELECT
		COUNT(*) FILTER (WHERE st.status = 'open')::int,
		COUNT(*) FILTER (WHERE st.priority IN ('high', 'urgent'))::int,
		COUNT(*) FILTER (WHERE u.role = 'business')::int,
		COUNT(*) FILTER (WHERE u.role = 'influencer')::int,
		COUNT(*) FILTER (WHERE st.status = 'resolved')::int
	FROM support_tickets st
	INNER JOIN users u ON u.id = st.created_by_user_id
	`

	getAdminSupportTicket = `
	SELECT
		st.id,
		st.subject,
		st.message,
		st.issue_type::text,
		st.priority::text,
		st.status::text,
		st.resolution_note,
		(
			SELECT sm.message
			FROM support_messages sm
			WHERE sm.ticket_id = st.id AND sm.is_internal_note = TRUE
			ORDER BY sm.created_at DESC
			LIMIT 1
		),
		u.id,
		COALESCE(bp.business_name, ip.creator_name, split_part(u.email, '@', 1)),
		u.email,
		u.role::text,
		st.related_campaign_id,
		c.title,
		st.created_at,
		st.updated_at
	FROM support_tickets st
	INNER JOIN users u ON u.id = st.created_by_user_id
	LEFT JOIN business_profiles bp ON bp.user_id = u.id AND u.role = 'business'
	LEFT JOIN influencer_profiles ip ON ip.user_id = u.id AND u.role = 'influencer'
	LEFT JOIN campaigns c ON c.id = st.related_campaign_id
	WHERE st.id = $1
	`

	updateSupportTicketStatus = `
	UPDATE support_tickets
	SET
		status = $2::support_ticket_status,
		updated_at = NOW(),
		resolved_at = CASE WHEN $2::text = 'resolved' THEN NOW() ELSE resolved_at END,
		resolved_by = CASE WHEN $2::text = 'resolved' THEN $3::uuid ELSE resolved_by END,
		assigned_admin_user_id = COALESCE(assigned_admin_user_id, $3::uuid)
	WHERE id = $1
	`

	updateSupportTicketReply = `
	UPDATE support_tickets
	SET
		resolution_note = $2,
		assigned_admin_user_id = $3,
		updated_at = NOW()
	WHERE id = $1
	`

	deleteInternalNotes = `
	DELETE FROM support_messages
	WHERE ticket_id = $1 AND is_internal_note = TRUE
	`

	insertInternalNote = `
	INSERT INTO support_messages (id, ticket_id, sender_user_id, message, is_internal_note)
	VALUES ($1, $2, $3, $4, TRUE)
	`

	listInfluencerInviteOptions = `
	SELECT ci.id, c.title, bp.business_name
	FROM campaign_invites ci
	INNER JOIN campaigns c ON c.id = ci.campaign_id
	INNER JOIN business_profiles bp ON bp.user_id = ci.business_user_id
	WHERE ci.influencer_user_id = $1
	ORDER BY ci.created_at DESC
	LIMIT 50
	`

	getInviteCampaignID = `
	SELECT ci.campaign_id
	FROM campaign_invites ci
	WHERE ci.id = $1 AND ci.influencer_user_id = $2
	`
)
