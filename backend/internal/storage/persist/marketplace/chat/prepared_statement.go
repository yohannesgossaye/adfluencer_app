package chat

// messageColumns is the column list every message query returns, in the
// order scanMessage expects.
const messageColumns = `
	id,
	conversation_id,
	sender_user_id,
	receiver_user_id,
	body,
	media_url,
	media_type,
	is_read,
	read_at,
	deleted_for_sender,
	deleted_for_receiver,
	is_deleted_for_everyone,
	client_message_id,
	created_at,
	updated_at`

// One conversation per participant pair (uq_conversations_pair, migration
// 00006), whatever campaign it was opened from; pair_low/pair_high hold the
// two user ids in a fixed order.
const (
	insertConversation = `
		INSERT INTO conversations (campaign_id, status, pair_low, pair_high, created_at, updated_at)
		VALUES ($1, 'active', LEAST($2::uuid, $3::uuid), GREATEST($2::uuid, $3::uuid), now(), now())
		ON CONFLICT (pair_low, pair_high) WHERE pair_low IS NOT NULL
		DO NOTHING
		RETURNING id
	`

	// The pair's existing chat; one opened without a campaign takes the
	// campaign it is now opened from, so it carries a useful label.
	findConversationByPair = `
		UPDATE conversations
		SET campaign_id = COALESCE(campaign_id, $3::uuid)
		WHERE pair_low = LEAST($1::uuid, $2::uuid)
		  AND pair_high = GREATEST($1::uuid, $2::uuid)
		RETURNING id
	`

	insertParticipant = `
		INSERT INTO conversation_participants (conversation_id, user_id, role, created_at)
		SELECT $1, u.id, u.role::text::conversation_participant_role, now()
		FROM users u WHERE u.id = $2
		ON CONFLICT (conversation_id, user_id) DO NOTHING
	`

	isBusinessVerified = `
		SELECT COALESCE((SELECT is_verified FROM business_profiles WHERE user_id = $1), FALSE)
	`

	getCampaignOwner = `
		SELECT business_user_id FROM campaigns WHERE id = $1
	`

	conversationSelect = `
		SELECT
			c.id,
			c.campaign_id,
			cam.title AS campaign_title,
			c.status,
			c.last_message_preview,
			c.last_message_at,
			c.created_at,
			u.id AS other_user_id,
			u.role AS other_user_role,
			COALESCE(bp.business_name, ip.creator_name, 'Adfluence user') AS other_user_name,
			COALESCE(bp.logo_url, ip.profile_image_url) AS other_user_avatar,
			(
				SELECT COUNT(*)
				FROM messages m
				WHERE m.conversation_id = c.id
				  AND m.receiver_user_id = $1
				  AND m.is_read = false
				  AND NOT m.is_deleted_for_everyone
				  AND NOT m.deleted_for_receiver
			) AS unread_count
		FROM conversations c
		JOIN conversation_participants my_p ON my_p.conversation_id = c.id AND my_p.user_id = $1
		JOIN conversation_participants other_p ON other_p.conversation_id = c.id AND other_p.user_id <> $1
		JOIN users u ON u.id = other_p.user_id
		LEFT JOIN business_profiles bp ON bp.user_id = u.id
		LEFT JOIN influencer_profiles ip ON ip.user_id = u.id
		LEFT JOIN campaigns cam ON cam.id = c.campaign_id
	`

	getConversationSummary = conversationSelect + `
		WHERE c.id = $2
	`

	listConversationsByUser = conversationSelect + `
		ORDER BY COALESCE(c.last_message_at, c.created_at) DESC, c.id DESC
		LIMIT $2 OFFSET $3
	`

	checkParticipantExists = `
		SELECT EXISTS (
			SELECT 1 FROM conversation_participants
			WHERE conversation_id = $1 AND user_id = $2
		)
	`

	getOtherParticipant = `
		SELECT user_id
		FROM conversation_participants
		WHERE conversation_id = $1 AND user_id <> $2
		LIMIT 1
	`

	// A repeated client_message_id from the same sender inserts nothing; the
	// caller then loads the original message.
	insertMessage = `
		INSERT INTO messages (
			conversation_id, sender_user_id, receiver_user_id, body, media_url, media_type,
			client_message_id, is_read, deleted_for_sender, deleted_for_receiver,
			is_deleted_for_everyone, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, false, false, false, false, now(), now())
		ON CONFLICT (sender_user_id, client_message_id) WHERE client_message_id IS NOT NULL
		DO NOTHING
		RETURNING ` + messageColumns

	getMessageByClientID = `
		SELECT ` + messageColumns + `
		FROM messages
		WHERE sender_user_id = $1 AND client_message_id = $2
	`

	updateConversationLastMessage = `
		UPDATE conversations
		SET last_message_preview = $1, last_message_at = now(), updated_at = now()
		WHERE id = $2
	`

	// Newest page first, then returned oldest→newest. (created_at, id) breaks
	// ties so messages sharing a timestamp are never skipped. $3 is the id of
	// the oldest message already shown (NULL for the newest page).
	getMessagesByConversation = `
		WITH page AS (
			SELECT
				id,
				conversation_id,
				sender_user_id,
				receiver_user_id,
				CASE WHEN is_deleted_for_everyone THEN 'This message was deleted' ELSE body END AS body,
				CASE WHEN is_deleted_for_everyone THEN NULL ELSE media_url END AS media_url,
				media_type,
				is_read,
				read_at,
				deleted_for_sender,
				deleted_for_receiver,
				is_deleted_for_everyone,
				client_message_id,
				created_at,
				updated_at
			FROM messages
			WHERE conversation_id = $1
			  AND (
				$3::uuid IS NULL
				OR (created_at, id) < (
					SELECT created_at, id FROM messages WHERE id = $3::uuid AND conversation_id = $1
				)
			  )
			  AND NOT (sender_user_id = $2 AND deleted_for_sender)
			  AND NOT (receiver_user_id = $2 AND deleted_for_receiver)
			ORDER BY created_at DESC, id DESC
			LIMIT $4
		)
		SELECT * FROM page ORDER BY created_at ASC, id ASC
	`

	markMessagesRead = `
		UPDATE messages
		SET is_read = true, read_at = now(), updated_at = now()
		WHERE conversation_id = $1 AND receiver_user_id = $2 AND is_read = false
	`

	updateParticipantLastRead = `
		UPDATE conversation_participants
		SET last_read_at = now()
		WHERE conversation_id = $1 AND user_id = $2
	`

	getMessageByID = `
		SELECT ` + messageColumns + `
		FROM messages
		WHERE id = $1
	`

	deleteMessageForEveryone = `
		UPDATE messages
		SET is_deleted_for_everyone = true, updated_at = now()
		WHERE id = $1
		RETURNING ` + messageColumns

	deleteMessageForSender = `
		UPDATE messages
		SET deleted_for_sender = true, updated_at = now()
		WHERE id = $1
		RETURNING ` + messageColumns

	deleteMessageForReceiver = `
		UPDATE messages
		SET deleted_for_receiver = true, updated_at = now()
		WHERE id = $1
		RETURNING ` + messageColumns

	// After a delete-for-everyone the list must not keep showing the text.
	refreshConversationPreview = `
		UPDATE conversations c
		SET last_message_preview = (
				SELECT CASE
					WHEN m.body = '' AND m.media_url IS NOT NULL THEN 'Attachment: ' || m.media_type
					WHEN char_length(m.body) > 100 THEN left(m.body, 97) || '...'
					ELSE m.body
				END
				FROM messages m
				WHERE m.conversation_id = c.id AND NOT m.is_deleted_for_everyone
				ORDER BY m.created_at DESC, m.id DESC
				LIMIT 1
			),
			updated_at = now()
		WHERE c.id = $1
	`
)
