-- +goose Up
-- One conversation per business–creator pair. Until now a chat was unique per
-- pair *and campaign*, so opening "Message" from a page without a campaign
-- (e.g. the creator's profile reached from a contract) started a second,
-- empty "Direct message" chat next to the real one. Merge every pair's chats
-- into its oldest one, then allow only one chat per pair.

CREATE TEMP TABLE chat_pair_dupes ON COMMIT DROP AS
SELECT id, keeper
FROM (
    SELECT id,
           first_value(id) OVER (PARTITION BY pair_low, pair_high ORDER BY created_at, id) AS keeper
    FROM conversations
    WHERE pair_low IS NOT NULL
) ranked
WHERE id <> keeper;

-- A keeper without a campaign takes the most recent campaign of its duplicates.
UPDATE conversations k
SET campaign_id = latest.campaign_id
FROM (
    SELECT DISTINCT ON (d.keeper) d.keeper, c.campaign_id
    FROM chat_pair_dupes d
    JOIN conversations c ON c.id = d.id
    WHERE c.campaign_id IS NOT NULL
    ORDER BY d.keeper, c.created_at DESC
) latest
WHERE k.id = latest.keeper AND k.campaign_id IS NULL;

UPDATE messages m SET conversation_id = d.keeper
FROM chat_pair_dupes d WHERE m.conversation_id = d.id;

UPDATE message_reports r SET conversation_id = d.keeper
FROM chat_pair_dupes d WHERE r.conversation_id = d.id;

UPDATE chat_admin_access_logs l SET conversation_id = d.keeper
FROM chat_pair_dupes d WHERE l.conversation_id = d.id;

UPDATE contracts k SET conversation_id = d.keeper
FROM chat_pair_dupes d WHERE k.conversation_id = d.id;

-- The preview and time of each merged chat follow its latest message.
UPDATE conversations c
SET last_message_at = lm.created_at,
    last_message_preview = CASE
        WHEN char_length(lm.body) > 100 THEN left(lm.body, 97) || '...'
        ELSE lm.body
    END
FROM (
    SELECT DISTINCT ON (conversation_id) conversation_id, body, created_at
    FROM messages
    WHERE NOT is_deleted_for_everyone
      AND conversation_id IN (SELECT keeper FROM chat_pair_dupes)
    ORDER BY conversation_id, created_at DESC, id DESC
) lm
WHERE c.id = lm.conversation_id;

-- The emptied duplicates go (their participant rows go with them).
DELETE FROM conversations c USING chat_pair_dupes d WHERE c.id = d.id;

DROP INDEX IF EXISTS uq_conversations_pair_campaign;
CREATE UNIQUE INDEX IF NOT EXISTS uq_conversations_pair
    ON conversations (pair_low, pair_high)
    WHERE pair_low IS NOT NULL;

-- +goose Down
-- Merged chats stay merged; only the old rule comes back.
DROP INDEX IF EXISTS uq_conversations_pair;
CREATE UNIQUE INDEX IF NOT EXISTS uq_conversations_pair_campaign
    ON conversations (pair_low, pair_high, (COALESCE(campaign_id, '00000000-0000-0000-0000-000000000000'::uuid)))
    WHERE pair_low IS NOT NULL;
