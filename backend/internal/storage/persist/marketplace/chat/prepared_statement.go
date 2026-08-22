package chat

const (
	findConversationWithCampaign = `
		SELECT c.id 
		FROM conversations c
		JOIN conversation_participants p1 ON p1.conversation_id = c.id AND p1.user_id = $1
		JOIN conversation_participants p2 ON p2.conversation_id = c.id AND p2.user_id = $2
		WHERE c.campaign_id = $3
		LIMIT 1
	`

	findConversationWithoutCampaign = `
		SELECT c.id 
		FROM conversations c
		JOIN conversation_participants p1 ON p1.conversation_id = c.id AND p1.user_id = $1
		JOIN conversation_participants p2 ON p2.conversation_id = c.id AND p2.user_id = $2
		LIMIT 1
	`

	getUserRole = `
		SELECT role FROM users WHERE id = $1
	`

	insertConversation = `
		INSERT INTO conversations (campaign_id, status, created_at, updated_at)
		VALUES ($1, 'active', now(), now())
		RETURNING id
	`

	insertParticipant = `
		INSERT INTO conversation_participants (conversation_id, user_id, role, created_at)
		VALUES ($1, $2, $3, now())
	`

	getConversationSummary = `
		SELECT 
			c.id,
			c.campaign_id,
			cam.title as campaign_title,
			c.status,
			c.last_message_preview,
			c.last_message_at,
			c.created_at,
			u.id as other_user_id,
			u.role as other_user_role,
			COALESCE(bp.business_name, ip.creator_name, u.email) as other_user_name,
			COALESCE(bp.logo_url, ip.profile_image_url) as other_user_avatar,
			(
				SELECT COUNT(*) 
				FROM messages m 
				WHERE m.conversation_id = c.id 
				  AND m.receiver_user_id = $2 
				  AND m.is_read = false
			) as unread_count
		FROM conversations c
		LEFT JOIN campaigns cam ON cam.id = c.campaign_id
		JOIN conversation_participants other_p ON other_p.conversation_id = c.id AND other_p.user_id <> $2
		JOIN users u ON u.id = other_p.user_id
		LEFT JOIN business_profiles bp ON bp.user_id = u.id
		LEFT JOIN influencer_profiles ip ON ip.user_id = u.id
		WHERE c.id = $1
		LIMIT 1
	`

	listConversationsByUser = `
		SELECT 
			c.id,
			c.campaign_id,
			cam.title as campaign_title,
			c.status,
			c.last_message_preview,
			c.last_message_at,
			c.created_at,
			u.id as other_user_id,
			u.role as other_user_role,
			COALESCE(bp.business_name, ip.creator_name, u.email) as other_user_name,
			COALESCE(bp.logo_url, ip.profile_image_url) as other_user_avatar,
			(
				SELECT COUNT(*) 
				FROM messages m 
				WHERE m.conversation_id = c.id 
				  AND m.receiver_user_id = $1 
				  AND m.is_read = false
			) as unread_count
		FROM conversations c
		JOIN conversation_participants my_p ON my_p.conversation_id = c.id AND my_p.user_id = $1
		JOIN conversation_participants other_p ON other_p.conversation_id = c.id AND other_p.user_id <> $1
		JOIN users u ON u.id = other_p.user_id
		LEFT JOIN business_profiles bp ON bp.user_id = u.id
		LEFT JOIN influencer_profiles ip ON ip.user_id = u.id
		LEFT JOIN campaigns cam ON cam.id = c.campaign_id
		ORDER BY COALESCE(c.last_message_at, c.created_at) DESC
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

	insertMessage = `
		INSERT INTO messages (
			conversation_id,
			sender_user_id,
			receiver_user_id,
			body,
			media_url,
			media_type,
			is_read,
			deleted_for_sender,
			deleted_for_receiver,
			is_deleted_for_everyone,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, false, false, false, false, now(), now())
		RETURNING 
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
			created_at,
			updated_at
	`

	updateConversationLastMessage = `
		UPDATE conversations 
		SET last_message_preview = $1, last_message_at = now(), updated_at = now()
		WHERE id = $2
	`

	getMessagesByConversation = `
		WITH raw_messages AS (
			SELECT 
				id,
				conversation_id,
				sender_user_id,
				receiver_user_id,
				CASE 
					WHEN is_deleted_for_everyone THEN 'This message was deleted'
					ELSE body 
				END as body,
				CASE 
					WHEN is_deleted_for_everyone THEN NULL
					ELSE media_url 
				END as media_url,
				media_type,
				is_read,
				read_at,
				deleted_for_sender,
				deleted_for_receiver,
				is_deleted_for_everyone,
				created_at,
				updated_at
			FROM messages
			WHERE conversation_id = $1
			  AND (
				$3 = '' 
				OR created_at < (SELECT created_at FROM messages WHERE id = $3::uuid)
			  )
			  AND NOT (sender_user_id = $2 AND deleted_for_sender = true)
			  AND NOT (receiver_user_id = $2 AND deleted_for_receiver = true)
			ORDER BY created_at DESC
			LIMIT $4
		)
		SELECT * FROM raw_messages ORDER BY created_at ASC
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
		SELECT 
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
			created_at,
			updated_at
		FROM messages
		WHERE id = $1
	`

	deleteMessageForEveryone = `
		UPDATE messages 
		SET is_deleted_for_everyone = true, updated_at = now()
		WHERE id = $1
		RETURNING 
			id, conversation_id, sender_user_id, receiver_user_id, body,
			media_url, media_type, is_read, read_at, deleted_for_sender,
			deleted_for_receiver, is_deleted_for_everyone, created_at, updated_at
	`

	deleteMessageForSender = `
		UPDATE messages 
		SET deleted_for_sender = true, updated_at = now()
		WHERE id = $1
		RETURNING 
			id, conversation_id, sender_user_id, receiver_user_id, body,
			media_url, media_type, is_read, read_at, deleted_for_sender,
			deleted_for_receiver, is_deleted_for_everyone, created_at, updated_at
	`

	deleteMessageForReceiver = `
		UPDATE messages 
		SET deleted_for_receiver = true, updated_at = now()
		WHERE id = $1
		RETURNING 
			id, conversation_id, sender_user_id, receiver_user_id, body,
			media_url, media_type, is_read, read_at, deleted_for_sender,
			deleted_for_receiver, is_deleted_for_everyone, created_at, updated_at
	`
)