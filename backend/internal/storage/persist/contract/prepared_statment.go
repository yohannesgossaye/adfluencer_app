package contract

const (
	getAcceptedInvite = `
	SELECT
		ci.id,
		ci.campaign_id,
		ci.business_user_id,
		ci.influencer_user_id,
		c.status::text
	FROM campaign_invites ci
	JOIN campaigns c ON c.id = ci.campaign_id
	WHERE ci.id = $1
		AND ci.business_user_id = $2
		AND ci.status = 'accepted'
	`

	contractExistsForInvite = `
	SELECT EXISTS (
		SELECT 1 FROM contracts WHERE invite_id = $1
	)
	`

	insertContract = `
	INSERT INTO contracts (
		id, campaign_id, business_user_id, influencer_user_id, invite_id,
		title, description, amount_etb,
		business_fee_percent, influencer_fee_percent,
		business_fee_amount_etb, influencer_fee_amount_etb,
		total_business_payable_etb, influencer_receivable_etb, platform_total_fee_etb,
		start_date, end_date, status,
		accepted_by_business_at
	) VALUES (
		$1, $2, $3, $4, $5,
		$6, $7, $8,
		$9, $10,
		$11, $12,
		$13, $14, $15,
		$16, $17, 'pending_influencer',
		NOW()
	)
	`

	insertContractDeliverable = `
	INSERT INTO contract_deliverables (
		id, contract_id, deliverable_type, quantity, description, status
	) VALUES ($1, $2, $3::deliverable_type, $4, $5, 'pending')
	`

	insertContractStatusHistory = `
	INSERT INTO contract_status_history (
		id, contract_id, old_status, new_status, changed_by, reason
	) VALUES ($1, $2, $3::contract_status, $4::contract_status, $5, $6)
	`

	getContractByID = `
	SELECT
		c.id,
		c.invite_id,
		c.campaign_id,
		c.business_user_id,
		c.influencer_user_id,
		c.title,
		c.description,
		c.amount_etb,
		c.business_fee_percent,
		c.influencer_fee_percent,
		c.business_fee_amount_etb,
		c.influencer_fee_amount_etb,
		c.total_business_payable_etb,
		c.influencer_receivable_etb,
		c.platform_total_fee_etb,
		c.start_date,
		c.end_date,
		c.status::text,
		c.accepted_by_influencer_at,
		c.completed_at,
		c.created_at,
		c.updated_at,
		camp.title,
		bp.business_name,
		ip.creator_name
	FROM contracts c
	JOIN campaigns camp ON camp.id = c.campaign_id
	JOIN business_profiles bp ON bp.user_id = c.business_user_id
	JOIN influencer_profiles ip ON ip.user_id = c.influencer_user_id
	WHERE c.id = $1
	`

	countBusinessContracts = `
	SELECT COUNT(*)::int
	FROM contracts c
	WHERE c.business_user_id = $1
		AND ($2 = '' OR c.status::text = $2)
	`

	listBusinessContracts = `
	SELECT
		c.id,
		c.title,
		c.status::text,
		c.amount_etb,
		camp.title,
		ip.creator_name,
		c.total_business_payable_etb,
		c.influencer_receivable_etb,
		c.created_at
	FROM contracts c
	JOIN campaigns camp ON camp.id = c.campaign_id
	JOIN influencer_profiles ip ON ip.user_id = c.influencer_user_id
	WHERE c.business_user_id = $1
		AND ($2 = '' OR c.status::text = $2)
	ORDER BY c.created_at DESC
	LIMIT $3 OFFSET $4
	`

	countInfluencerContracts = `
	SELECT COUNT(*)::int
	FROM contracts c
	WHERE c.influencer_user_id = $1
		AND ($2 = '' OR c.status::text = $2)
	`

	listInfluencerContracts = `
	SELECT
		c.id,
		c.title,
		c.status::text,
		c.amount_etb,
		camp.title,
		bp.business_name,
		c.total_business_payable_etb,
		c.influencer_receivable_etb,
		c.created_at
	FROM contracts c
	JOIN campaigns camp ON camp.id = c.campaign_id
	JOIN business_profiles bp ON bp.user_id = c.business_user_id
	WHERE c.influencer_user_id = $1
		AND ($2 = '' OR c.status::text = $2)
	ORDER BY c.created_at DESC
	LIMIT $3 OFFSET $4
	`

	acceptContract = `
	UPDATE contracts
	SET
		status = 'accepted',
		accepted_by_influencer_at = NOW(),
		updated_at = NOW()
	WHERE id = $1
		AND influencer_user_id = $2
		AND status = 'pending_influencer'
	`

	rejectContract = `
	UPDATE contracts
	SET status = 'rejected', updated_at = NOW()
	WHERE id = $1
		AND influencer_user_id = $2
		AND status = 'pending_influencer'
	`

	cancelBusinessContract = `
	UPDATE contracts
	SET status = 'cancelled', cancelled_at = NOW(), updated_at = NOW()
	WHERE id = $1
		AND business_user_id = $2
		AND status IN ('pending_influencer', 'accepted')
	`

	paidPaymentExists = `
	SELECT EXISTS (
		SELECT 1 FROM payments
		WHERE contract_id = $1 AND status = 'paid'
	)
	`

	insertPayment = `
	INSERT INTO payments (
		id, contract_id, business_user_id, influencer_user_id,
		amount_etb, business_fee_etb, influencer_fee_etb, platform_fee_total_etb,
		total_paid_by_business_etb, influencer_payout_etb,
		payment_method, status, paid_at
	) VALUES (
		$1, $2, $3, $4,
		$5, $6, $7, $8,
		$9, $10,
		$11, 'paid', NOW()
	)
	`

	activateContractAfterPay = `
	UPDATE contracts
	SET status = 'active', updated_at = NOW()
	WHERE id = $1 AND business_user_id = $2 AND status = 'accepted'
	`

	insertPayout = `
	INSERT INTO payouts (
		id, payment_id, influencer_user_id, amount_etb, status
	) VALUES ($1, $2, $3, $4, 'pending')
	`

	completeContractAfterApprove = `
	UPDATE contracts
	SET status = 'completed', completed_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND business_user_id = $2 AND status = 'active'
	`

	getContractPayment = `
	SELECT id, status::text, total_paid_by_business_etb, payment_method, paid_at
	FROM payments
	WHERE contract_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	getContractPayout = `
	SELECT p.id, p.status::text, p.amount_etb, p.created_at, p.paid_at
	FROM payouts p
	JOIN payments pay ON pay.id = p.payment_id
	WHERE pay.contract_id = $1
	ORDER BY p.created_at DESC
	LIMIT 1
	`

	getContractDeliverables = `
	SELECT id, deliverable_type::text, quantity, description, status::text
	FROM contract_deliverables
	WHERE contract_id = $1
	ORDER BY created_at
	`
)
