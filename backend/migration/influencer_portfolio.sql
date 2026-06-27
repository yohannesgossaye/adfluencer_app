-- Content samples (images / short videos) for influencer profiles — max 3 per user.

CREATE TABLE IF NOT EXISTS influencer_portfolio_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    file_url TEXT NOT NULL,
    object_key TEXT NOT NULL,

    media_type VARCHAR(10) NOT NULL
        CHECK (media_type IN ('image', 'video')),

    mime_type VARCHAR(100) NOT NULL,
    file_size_bytes BIGINT NOT NULL
        CHECK (file_size_bytes > 0),

    sort_order SMALLINT NOT NULL
        CHECK (sort_order >= 0 AND sort_order <= 2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (user_id, sort_order)
);

CREATE INDEX IF NOT EXISTS idx_influencer_portfolio_items_user_id
    ON influencer_portfolio_items(user_id);
