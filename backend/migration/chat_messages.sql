-- =============================================================================
-- CHAT_MESSAGES.SQL — Realtime chat messages table & indexes
-- =============================================================================
-- messages — 1-on-1 messages with read receipts, attachments, and deletion
-- =============================================================================

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,

    sender_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    receiver_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    body TEXT NOT NULL,

    media_url TEXT,
    media_type VARCHAR(50) DEFAULT 'text',

    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ,

    deleted_for_sender BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_for_receiver BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted_for_everyone BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_created
    ON messages(conversation_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_messages_sender
    ON messages(sender_user_id);

CREATE INDEX IF NOT EXISTS idx_messages_receiver
    ON messages(receiver_user_id);

CREATE INDEX IF NOT EXISTS idx_messages_is_read
    ON messages(conversation_id, receiver_user_id, is_read);
