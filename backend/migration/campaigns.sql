-- =============================================================================
-- CAMPAIGNS.SQL — Business campaigns + influencer matching (5 tables)
-- =============================================================================
-- campaigns                   — business campaign brief
-- campaign_platforms            — tiktok, instagram, youtube…
-- campaign_deliverables         — what content is needed
-- campaign_recommendations      — ranked match results (recommendation page)
-- campaign_invites              — business invites influencer
-- =============================================================================

CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    business_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    goal TEXT,
    description TEXT NOT NULL,

    category_id UUID REFERENCES categories(id),

    target_audience TEXT,
    city VARCHAR(100),
    region VARCHAR(100),

    budget_min_etb NUMERIC(14,2),
    budget_max_etb NUMERIC(14,2),

    creator_size_preference creator_size NOT NULL DEFAULT 'any',
    number_of_creators_needed INTEGER NOT NULL DEFAULT 1,

    start_date DATE,
    end_date DATE,

    status campaign_status NOT NULL DEFAULT 'draft',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE campaign_platforms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID NOT NULL
        REFERENCES campaigns(id) ON DELETE CASCADE,

    platform social_platform NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (campaign_id, platform)
);

CREATE TABLE campaign_deliverables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID NOT NULL
        REFERENCES campaigns(id) ON DELETE CASCADE,

    deliverable_type deliverable_type NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE campaign_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID NOT NULL
        REFERENCES campaigns(id) ON DELETE CASCADE,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    match_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    score_breakdown JSONB NOT NULL DEFAULT '{}',
    rank_position INTEGER,

    status recommendation_status NOT NULL DEFAULT 'recommended',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (campaign_id, influencer_user_id)
);

CREATE TABLE campaign_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID NOT NULL
        REFERENCES campaigns(id) ON DELETE CASCADE,
    business_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    message TEXT,

    status invite_status NOT NULL DEFAULT 'pending',

    responded_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (campaign_id, influencer_user_id)
);

CREATE INDEX idx_campaigns_business_user_id ON campaigns(business_user_id);
CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_campaigns_category_id ON campaigns(category_id);
CREATE INDEX idx_campaigns_created_at ON campaigns(created_at);

CREATE INDEX idx_campaign_platforms_campaign_id ON campaign_platforms(campaign_id);

CREATE INDEX idx_campaign_deliverables_campaign_id ON campaign_deliverables(campaign_id);

CREATE INDEX idx_campaign_recommendations_campaign_id ON campaign_recommendations(campaign_id);
CREATE INDEX idx_campaign_recommendations_influencer ON campaign_recommendations(influencer_user_id);
CREATE INDEX idx_campaign_recommendations_status ON campaign_recommendations(status);
CREATE INDEX idx_campaign_recommendations_rank ON campaign_recommendations(campaign_id, rank_position);

CREATE INDEX idx_campaign_invites_campaign_id ON campaign_invites(campaign_id);
CREATE INDEX idx_campaign_invites_business_user_id ON campaign_invites(business_user_id);
CREATE INDEX idx_campaign_invites_influencer_user_id ON campaign_invites(influencer_user_id);
CREATE INDEX idx_campaign_invites_status ON campaign_invites(status);
