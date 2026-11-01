-- Chat integrity (Phase 1B). Safe to re-run. Runs in one transaction.
--
-- 1. conversations.pair_low / pair_high: the two participants in a fixed
--    order, so "one conversation per pair (per campaign)" can be enforced by a
--    unique index instead of racy check-then-insert code.
-- 2. Existing duplicate conversations for the same pair+campaign are merged
--    into the oldest one: messages and references move over, duplicates are
--    deleted. No message is lost.
-- 3. messages.client_message_id: lets clients retry a send without creating
--    a duplicate message.
-- 4. Indexes for history paging and unread counts.

BEGIN;

ALTER TABLE conversations ADD COLUMN IF NOT EXISTS pair_low UUID;
ALTER TABLE conversations ADD COLUMN IF NOT EXISTS pair_high UUID;
ALTER TABLE messages ADD COLUMN IF NOT EXISTS client_message_id UUID;

UPDATE conversations c
SET pair_low = p1.user_id, pair_high = p2.user_id
FROM conversation_participants p1
JOIN conversation_participants p2
  ON p2.conversation_id = p1.conversation_id AND p1.user_id < p2.user_id
WHERE p1.conversation_id = c.id
  AND c.pair_low IS NULL;

CREATE TEMP TABLE chat_conversation_dupes ON COMMIT DROP AS
SELECT id, keeper
FROM (
    SELECT id,
           first_value(id) OVER (
               PARTITION BY pair_low, pair_high, campaign_id
               ORDER BY created_at, id
           ) AS keeper
    FROM conversations
    WHERE pair_low IS NOT NULL
) ranked
WHERE id <> keeper;

UPDATE messages m SET conversation_id = d.keeper
FROM chat_conversation_dupes d WHERE m.conversation_id = d.id;

UPDATE message_reports r SET conversation_id = d.keeper
FROM chat_conversation_dupes d WHERE r.conversation_id = d.id;

UPDATE chat_admin_access_logs l SET conversation_id = d.keeper
FROM chat_conversation_dupes d WHERE l.conversation_id = d.id;

UPDATE contracts k SET conversation_id = d.keeper
FROM chat_conversation_dupes d WHERE k.conversation_id = d.id;

-- Refresh the preview of every conversation that absorbed messages.
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
      AND conversation_id IN (SELECT keeper FROM chat_conversation_dupes)
    ORDER BY conversation_id, created_at DESC, id DESC
) lm
WHERE c.id = lm.conversation_id;

DELETE FROM conversations c
USING chat_conversation_dupes d
WHERE c.id = d.id;

CREATE UNIQUE INDEX IF NOT EXISTS uq_conversations_pair_campaign
    ON conversations (pair_low, pair_high, (COALESCE(campaign_id, '00000000-0000-0000-0000-000000000000'::uuid)))
    WHERE pair_low IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_messages_sender_client_message_id
    ON messages (sender_user_id, client_message_id)
    WHERE client_message_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_messages_conversation_created_id
    ON messages (conversation_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_messages_unread
    ON messages (receiver_user_id, conversation_id)
    WHERE is_read = false;

COMMIT;
