package admin

const (
	listAdminDisputes = `
	SELECT
		d.id,
		d.description,
		d.issue_type::text,
		d.severity::text,
		d.status::text,
		d.created_at,
		camp.title,
		bp.business_name,
		ip.creator_name,
		COUNT(*) OVER()::int AS total_count
	FROM disputes d
	INNER JOIN contracts ct ON ct.id = d.contract_id
	INNER JOIN campaigns camp ON camp.id = ct.campaign_id
	INNER JOIN business_profiles bp ON bp.user_id = ct.business_user_id
	INNER JOIN influencer_profiles ip ON ip.user_id = ct.influencer_user_id
	WHERE (
			($1 = 'open' AND d.status IN ('open', 'under_review', 'escalated'))
			OR ($1 = 'resolved' AND d.status IN ('resolved', 'closed'))
		)
		AND (
			$2 = ''
			OR d.description ILIKE '%' || $2 || '%'
			OR camp.title ILIKE '%' || $2 || '%'
			OR bp.business_name ILIKE '%' || $2 || '%'
			OR ip.creator_name ILIKE '%' || $2 || '%'
		)
	ORDER BY d.created_at DESC
	LIMIT $3 OFFSET $4
	`

	getAdminDispute = `
	SELECT
		d.id,
		d.contract_id,
		camp.id,
		camp.title,
		bp.business_name,
		ip.creator_name,
		u_inf.email,
		ct.influencer_user_id,
		d.against_user_id,
		d.issue_type::text,
		d.severity::text,
		d.status::text,
		d.description,
		d.resolution_note,
		d.created_at,
		d.updated_at,
		d.resolved_at,
		COALESCE((
			SELECT al.description
			FROM audit_logs al
			WHERE al.entity_type = 'dispute_note'
				AND al.entity_id = d.id
			ORDER BY al.created_at DESC
			LIMIT 1
		), '')
	FROM disputes d
	INNER JOIN contracts ct ON ct.id = d.contract_id
	INNER JOIN campaigns camp ON camp.id = ct.campaign_id
	INNER JOIN business_profiles bp ON bp.user_id = ct.business_user_id
	INNER JOIN influencer_profiles ip ON ip.user_id = ct.influencer_user_id
	INNER JOIN users u_inf ON u_inf.id = ct.influencer_user_id
	WHERE d.id = $1
	`

	updateDisputeStatus = `
	UPDATE disputes
	SET
		status = $2::dispute_status,
		resolved_at = CASE WHEN $2::text IN ('resolved', 'closed') THEN NOW() ELSE NULL END,
		resolved_by = CASE WHEN $2::text IN ('resolved', 'closed') THEN $3::uuid ELSE NULL END,
		updated_at = NOW()
	WHERE id = $1
	`

	updateDisputeResolutionNote = `
	UPDATE disputes
	SET resolution_note = $2, updated_at = NOW()
	WHERE id = $1
	`

	getDisputeNotes = `
	SELECT
		d.resolution_note,
		COALESCE((
			SELECT al.description
			FROM audit_logs al
			WHERE al.entity_type = 'dispute_note'
				AND al.entity_id = d.id
			ORDER BY al.created_at DESC
			LIMIT 1
		), '')
	FROM disputes d
	WHERE d.id = $1
	`
)
