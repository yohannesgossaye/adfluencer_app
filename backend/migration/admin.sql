-- =============================================================================
-- ADMIN.SQL — Admin-managed reference data (3 tables)
-- =============================================================================
-- categories                    — niches/industries (beauty, fashion, tech…)
-- influencer_profile_categories — influencer ↔ category (normalized niches)
-- influencer_profile_platforms  — influencer ↔ platform (normalized platforms)
-- =============================================================================

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    description TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE influencer_profile_categories (
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL
        REFERENCES categories(id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (influencer_user_id, category_id)
);

CREATE TABLE influencer_profile_platforms (
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    platform social_platform NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (influencer_user_id, platform)
);

CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_is_active ON categories(is_active);
CREATE INDEX idx_influencer_profile_categories_category ON influencer_profile_categories(category_id);
CREATE INDEX idx_influencer_profile_platforms_platform ON influencer_profile_platforms(platform);
