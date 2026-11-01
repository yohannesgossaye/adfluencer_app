-- =============================================================================
-- AUTH.SQL — Login, profiles, verification, sessions (10 tables)
-- =============================================================================
-- users
-- business_profiles
-- influencer_profiles
-- admin_profiles
-- verification_documents
-- bank_accounts
-- user_sessions
-- password_reset_tokens
-- email_verification_tokens
-- audit_logs
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums (all types — must run before any table)
CREATE TYPE user_role AS ENUM ('business', 'influencer', 'admin');

CREATE TYPE user_status AS ENUM (
    'pending_profile',
    'pending_approval',
    'active',
    'rejected',
    'suspended',
    'blocked',
    'removed'
);

CREATE TYPE document_type AS ENUM (
    'fayda_id',
    'passport',
    'business_license',
    'trade_license',
    'tax_certificate',
    'portfolio_sample',
    'other'
);

CREATE TYPE document_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TYPE audit_action AS ENUM (
    'login',
    'logout',
    'signup',
    'password_reset',
    'profile_update',
    'status_change',
    'approve',
    'reject',
    'block',
    'unblock',
    'suspend',
    'unsuspend',
    'delete',
    'contract_status_change',
    'payment_status_change',
    'dispute_resolved',
    'permission_update',
    'admin_action'
);

CREATE TYPE campaign_status AS ENUM (
    'draft', 'pending', 'active', 'completed', 'cancelled', 'paused'
);

CREATE TYPE social_platform AS ENUM (
    'tiktok', 'instagram', 'youtube', 'telegram', 'facebook'
);

CREATE TYPE deliverable_type AS ENUM (
    'tiktok_video',
    'instagram_reel',
    'instagram_story',
    'youtube_short',
    'product_review',
    'post',
    'other'
);

CREATE TYPE creator_size AS ENUM ('nano', 'micro', 'mid', 'macro', 'mega', 'any');

CREATE TYPE recommendation_status AS ENUM ('recommended', 'invited', 'dismissed');

CREATE TYPE invite_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled');

CREATE TYPE conversation_status AS ENUM ('active', 'archived', 'closed');

CREATE TYPE conversation_participant_role AS ENUM ('business', 'influencer');

CREATE TYPE message_report_status AS ENUM ('open', 'under_review', 'resolved', 'dismissed');

CREATE TYPE contract_status AS ENUM (
    'draft',
    'pending_influencer',
    'pending_business',
    'accepted',
    'active',
    'completed',
    'cancelled',
    'disputed',
    'rejected'
);

CREATE TYPE contract_deliverable_status AS ENUM (
    'pending', 'submitted', 'approved', 'revision_requested', 'completed'
);

CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded', 'cancelled');

CREATE TYPE payout_status AS ENUM ('pending', 'processing', 'paid', 'failed', 'cancelled');

CREATE TYPE support_issue_type AS ENUM (
    'campaign_issue',
    'contract_issue',
    'payment_question',
    'creator_communication',
    'profile_approval',
    'technical_issue',
    'general_feedback',
    'complaint'
);

CREATE TYPE support_ticket_status AS ENUM (
    'open', 'in_review', 'resolved', 'escalated', 'closed'
);

CREATE TYPE support_priority AS ENUM ('low', 'normal', 'high', 'urgent');

CREATE TYPE dispute_status AS ENUM ('open', 'under_review', 'resolved', 'escalated', 'closed');

CREATE TYPE dispute_severity AS ENUM ('low', 'medium', 'high');

CREATE TYPE notification_type AS ENUM (
    'campaign_invite',
    'invite_accepted',
    'invite_rejected',
    'contract_created',
    'contract_accepted',
    'contract_completed',
    'payment_paid',
    'payout_paid',
    'support_reply',
    'verification_approved',
    'verification_rejected',
    'profile_approved',
    'profile_rejected',
    'dispute_opened',
    'dispute_resolved',
    'message_report',
    'general'
);

-- Tables
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(50) UNIQUE,

    password_hash TEXT NOT NULL,

    role user_role NOT NULL,
    status user_status NOT NULL DEFAULT 'pending_profile',

    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,

    last_login_at TIMESTAMPTZ,

    suspended_until TIMESTAMPTZ,
    status_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE business_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE
        REFERENCES users(id) ON DELETE CASCADE,

    owner_full_name VARCHAR(255) NOT NULL,

    business_name VARCHAR(255) NOT NULL,
    legal_business_name VARCHAR(255),

    logo_url TEXT,
    description TEXT,

    industry VARCHAR(100),
    business_type VARCHAR(100),
    business_size VARCHAR(100),

    city VARCHAR(100),
    region VARCHAR(100),
    address TEXT,

    tin_number VARCHAR(100),

    website_url TEXT,
    instagram_url TEXT,
    tiktok_url TEXT,
    facebook_url TEXT,
    linkedin_url TEXT,

    preferred_platforms TEXT[],

    monthly_marketing_budget NUMERIC(14,2),

    is_verified BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE influencer_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE
        REFERENCES users(id) ON DELETE CASCADE,

    full_name VARCHAR(255) NOT NULL,
    creator_name VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE,

    profile_image_url TEXT,
    bio TEXT,

    city VARCHAR(100),
    region VARCHAR(100),

    languages TEXT[],

    social_links JSONB NOT NULL DEFAULT '{}',

    instagram_followers INTEGER NOT NULL DEFAULT 0,
    tiktok_followers INTEGER NOT NULL DEFAULT 0,
    youtube_subscribers INTEGER NOT NULL DEFAULT 0,
    total_followers INTEGER NOT NULL DEFAULT 0,

    engagement_rate NUMERIC(5,2),

    audience_description TEXT,
    content_style TEXT,
    portfolio_description TEXT,

    min_price_etb NUMERIC(12,2),
    max_price_etb NUMERIC(12,2),

    approval_note TEXT,

    is_verified BOOLEAN NOT NULL DEFAULT FALSE,

    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admin_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE
        REFERENCES users(id) ON DELETE CASCADE,

    full_name VARCHAR(255) NOT NULL,

    permission_level VARCHAR(50) NOT NULL DEFAULT 'super_admin',

    permissions JSONB NOT NULL DEFAULT '{
        "manage_users": true,
        "manage_campaigns": true,
        "manage_payments": true,
        "manage_disputes": true,
        "manage_categories": true,
        "manage_support": true,
        "view_reports": true
    }',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE verification_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    document_type document_type NOT NULL,

    file_url TEXT NOT NULL,
    file_name VARCHAR(255),
    mime_type VARCHAR(100),
    file_size_bytes BIGINT,

    status document_status NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,

    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    bank_name VARCHAR(150),
    account_holder_name VARCHAR(255),
    account_number VARCHAR(100),

    mobile_money_provider VARCHAR(100),
    mobile_money_number VARCHAR(50),

    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    refresh_token_hash TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,

    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE phone_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    actor_user_id UUID REFERENCES users(id),
    target_user_id UUID REFERENCES users(id),

    action audit_action NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,

    description TEXT,
    metadata JSONB,

    ip_address INET,
    user_agent TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_business_profiles_user_id ON business_profiles(user_id);
CREATE INDEX idx_influencer_profiles_user_id ON influencer_profiles(user_id);
CREATE INDEX idx_influencer_profiles_is_verified ON influencer_profiles(is_verified);
CREATE INDEX idx_influencer_profiles_username ON influencer_profiles(username);

CREATE INDEX idx_verification_documents_user_id ON verification_documents(user_id);
CREATE INDEX idx_verification_documents_status ON verification_documents(status);

CREATE INDEX idx_bank_accounts_user_id ON bank_accounts(user_id);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);
CREATE INDEX idx_phone_verification_tokens_user_id ON phone_verification_tokens(user_id);

CREATE INDEX idx_audit_logs_actor ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_target ON audit_logs(target_user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
-- Auth security: login lockouts + rate-limit event log
-- Safe to re-run (IF NOT EXISTS). For existing DB: make migrate-security

CREATE TABLE IF NOT EXISTS auth_login_lockouts (
    email VARCHAR(255) PRIMARY KEY,

    failed_attempts INT NOT NULL DEFAULT 0,
    lockout_tier INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth_rate_limit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    scope VARCHAR(50) NOT NULL,
    rate_key VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_auth_rate_limit_scope_key_created
    ON auth_rate_limit_events (scope, rate_key, created_at DESC);
-- Phone OTP verification tokens (mirrors email_verification_tokens).
-- Safe to re-run (IF NOT EXISTS). For existing DB: make migrate-phone-verification

CREATE TABLE IF NOT EXISTS phone_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_phone_verification_tokens_user_id ON phone_verification_tokens(user_id);
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
-- =============================================================================
-- MARKETPLACE.SQL — Chat, contracts, payments (9 tables)
-- =============================================================================
-- conversations
-- conversation_participants
-- message_reports
-- chat_admin_access_logs
-- contracts
-- contract_deliverables
-- contract_status_history
-- payments
-- payouts
-- =============================================================================
-- Firebase stores message bodies: conversations/{id}/messages/{messageId}
-- Link chat to contract: contracts.conversation_id -> conversations.id
-- =============================================================================

CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,

    status conversation_status NOT NULL DEFAULT 'active',

    last_message_preview TEXT,
    last_message_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    role conversation_participant_role NOT NULL,

    last_read_at TIMESTAMPTZ,
    is_muted BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (conversation_id, user_id)
);

CREATE TABLE message_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,
    firebase_message_id VARCHAR(255) NOT NULL,

    reporter_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    reported_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    reason TEXT NOT NULL,

    status message_report_status NOT NULL DEFAULT 'open',

    resolved_by UUID REFERENCES users(id),
    resolution_note TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE chat_admin_access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    admin_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,

    reason TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID NOT NULL
        REFERENCES campaigns(id) ON DELETE RESTRICT,
    business_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,

    conversation_id UUID REFERENCES conversations(id) ON DELETE SET NULL,
    invite_id UUID REFERENCES campaign_invites(id) ON DELETE SET NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,

    amount_etb NUMERIC(14,2) NOT NULL,

    business_fee_percent NUMERIC(5,2) NOT NULL DEFAULT 5.00,
    influencer_fee_percent NUMERIC(5,2) NOT NULL DEFAULT 5.00,

    business_fee_amount_etb NUMERIC(14,2) NOT NULL,
    influencer_fee_amount_etb NUMERIC(14,2) NOT NULL,
    total_business_payable_etb NUMERIC(14,2) NOT NULL,
    influencer_receivable_etb NUMERIC(14,2) NOT NULL,
    platform_total_fee_etb NUMERIC(14,2) NOT NULL,

    start_date DATE,
    end_date DATE,

    status contract_status NOT NULL DEFAULT 'draft',

    accepted_by_business_at TIMESTAMPTZ,
    accepted_by_influencer_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contract_deliverables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE CASCADE,

    deliverable_type deliverable_type NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    description TEXT,

    status contract_deliverable_status NOT NULL DEFAULT 'pending',

    due_date DATE,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contract_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE CASCADE,

    old_status contract_status,
    new_status contract_status NOT NULL,

    changed_by UUID REFERENCES users(id),
    reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE RESTRICT,
    business_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,

    amount_etb NUMERIC(14,2) NOT NULL,
    business_fee_etb NUMERIC(14,2) NOT NULL,
    influencer_fee_etb NUMERIC(14,2) NOT NULL,
    platform_fee_total_etb NUMERIC(14,2) NOT NULL,
    total_paid_by_business_etb NUMERIC(14,2) NOT NULL,
    influencer_payout_etb NUMERIC(14,2) NOT NULL,

    payment_method VARCHAR(50) NOT NULL,
    provider VARCHAR(100),
    provider_reference VARCHAR(255),

    status payment_status NOT NULL DEFAULT 'pending',

    paid_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID NOT NULL
        REFERENCES payments(id) ON DELETE RESTRICT,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    bank_account_id UUID
        REFERENCES bank_accounts(id) ON DELETE SET NULL,

    amount_etb NUMERIC(14,2) NOT NULL,

    status payout_status NOT NULL DEFAULT 'pending',

    provider VARCHAR(100),
    provider_reference VARCHAR(255),

    paid_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_conversations_campaign_id ON conversations(campaign_id);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC);

CREATE INDEX idx_conversation_participants_conversation ON conversation_participants(conversation_id);
CREATE INDEX idx_conversation_participants_user ON conversation_participants(user_id);

CREATE INDEX idx_message_reports_conversation ON message_reports(conversation_id);
CREATE INDEX idx_message_reports_status ON message_reports(status);
CREATE INDEX idx_message_reports_reporter ON message_reports(reporter_user_id);

CREATE INDEX idx_chat_admin_access_logs_admin ON chat_admin_access_logs(admin_user_id);
CREATE INDEX idx_chat_admin_access_logs_conversation ON chat_admin_access_logs(conversation_id);

CREATE INDEX idx_contracts_campaign_id ON contracts(campaign_id);
CREATE INDEX idx_contracts_business_user_id ON contracts(business_user_id);
CREATE INDEX idx_contracts_influencer_user_id ON contracts(influencer_user_id);
CREATE INDEX idx_contracts_conversation_id ON contracts(conversation_id);
CREATE INDEX idx_contracts_status ON contracts(status);
CREATE INDEX idx_contracts_created_at ON contracts(created_at);

CREATE INDEX idx_contract_deliverables_contract_id ON contract_deliverables(contract_id);
CREATE INDEX idx_contract_deliverables_status ON contract_deliverables(status);

CREATE INDEX idx_contract_status_history_contract_id ON contract_status_history(contract_id);
CREATE INDEX idx_contract_status_history_created_at ON contract_status_history(created_at);

CREATE INDEX idx_payments_contract_id ON payments(contract_id);
CREATE INDEX idx_payments_business_user_id ON payments(business_user_id);
CREATE INDEX idx_payments_influencer_user_id ON payments(influencer_user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_paid_at ON payments(paid_at);

CREATE INDEX idx_payouts_payment_id ON payouts(payment_id);
CREATE INDEX idx_payouts_influencer_user_id ON payouts(influencer_user_id);
CREATE INDEX idx_payouts_status ON payouts(status);
CREATE INDEX idx_payouts_paid_at ON payouts(paid_at);
-- =============================================================================
-- SUPPORT.SQL — Support, disputes, reviews, notifications (5 tables)
-- =============================================================================
-- support_tickets
-- support_messages
-- disputes
-- reviews
-- notifications
-- =============================================================================

CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    created_by_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    assigned_admin_user_id UUID
        REFERENCES users(id) ON DELETE SET NULL,

    related_campaign_id UUID
        REFERENCES campaigns(id) ON DELETE SET NULL,
    related_contract_id UUID
        REFERENCES contracts(id) ON DELETE SET NULL,

    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,

    issue_type support_issue_type NOT NULL,
    priority support_priority NOT NULL DEFAULT 'normal',
    status support_ticket_status NOT NULL DEFAULT 'open',

    resolution_note TEXT,
    resolved_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE support_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    ticket_id UUID NOT NULL
        REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    message TEXT NOT NULL,
    attachment_url TEXT,

    is_internal_note BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE RESTRICT,

    opened_by_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    against_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    assigned_admin_user_id UUID
        REFERENCES users(id) ON DELETE SET NULL,

    issue_type support_issue_type NOT NULL,
    description TEXT NOT NULL,

    severity dispute_severity NOT NULL DEFAULT 'medium',
    status dispute_status NOT NULL DEFAULT 'open',

    resolution_note TEXT,
    resolved_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE CASCADE,

    reviewer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    reviewed_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (contract_id, reviewer_user_id, reviewed_user_id)
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,

    type notification_type NOT NULL,

    entity_type VARCHAR(100),
    entity_id UUID,

    read_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_support_tickets_created_by ON support_tickets(created_by_user_id);
CREATE INDEX idx_support_tickets_assigned_admin ON support_tickets(assigned_admin_user_id);
CREATE INDEX idx_support_tickets_status ON support_tickets(status);
CREATE INDEX idx_support_tickets_issue_type ON support_tickets(issue_type);
CREATE INDEX idx_support_tickets_related_campaign ON support_tickets(related_campaign_id);
CREATE INDEX idx_support_tickets_related_contract ON support_tickets(related_contract_id);

CREATE INDEX idx_support_messages_ticket_id ON support_messages(ticket_id);
CREATE INDEX idx_support_messages_sender ON support_messages(sender_user_id);

CREATE INDEX idx_disputes_contract_id ON disputes(contract_id);
CREATE INDEX idx_disputes_opened_by ON disputes(opened_by_user_id);
CREATE INDEX idx_disputes_against_user ON disputes(against_user_id);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_severity ON disputes(severity);

CREATE INDEX idx_reviews_contract_id ON reviews(contract_id);
CREATE INDEX idx_reviews_reviewer ON reviews(reviewer_user_id);
CREATE INDEX idx_reviews_reviewed_user ON reviews(reviewed_user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_read_at ON notifications(user_id, read_at);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
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
-- Seed categories for influencer niches (matches frontend profile.tsx).
-- Safe to re-run. For existing DB: make seed-categories

INSERT INTO categories (id, name, slug, description, is_active, created_at, updated_at)
VALUES
    ('c1000000-0000-4000-8000-000000000001', 'Beauty', 'beauty', 'Beauty and skincare creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000002', 'Fashion', 'fashion', 'Fashion and styling creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000003', 'Food', 'food', 'Food and restaurant creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000004', 'Tech', 'tech', 'Technology and gadget reviewers', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000005', 'Fitness', 'fitness', 'Fitness and wellness creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000006', 'Lifestyle', 'lifestyle', 'General lifestyle creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000007', 'Travel', 'travel', 'Travel and tourism creators', TRUE, NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;

-- Session rotation hardening (Phase 1A). Safe to re-run.
--
-- refresh_ttl_seconds: lifetime chosen at login (e.g. 30 days with "remember me"),
--   carried over every time the refresh token is rotated.
-- rotated_at: when this refresh token was exchanged for a new one. A rotated
--   token stays usable for a short grace window so that two browser tabs
--   refreshing at the same moment do not log each other out.

ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS refresh_ttl_seconds INTEGER;
ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS rotated_at TIMESTAMPTZ;

-- Every refresh and logout looks sessions up by token hash.
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_sessions_refresh_token_hash
    ON user_sessions (refresh_token_hash);

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

-- Phase 2A business logic. Safe to re-run. Runs in one transaction.
--
-- 1. Business verification becomes an explicit workflow the business can see:
--    unsubmitted -> pending -> verified | rejected (with a reason) -> pending ...
--    is_verified stays in sync (TRUE only when verification_status = 'verified').
-- 2. users.previous_status remembers the status before an admin block, so
--    unblocking restores it instead of making everyone "active".
-- 3. Campaigns paused by the old admin "pause" action were stored as
--    'pending'; they are moved to 'paused', the status pause now uses.

BEGIN;

ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verification_status TEXT NOT NULL DEFAULT 'unsubmitted';
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verification_note TEXT;
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verification_submitted_at TIMESTAMPTZ;
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ;
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verified_by UUID REFERENCES users(id) ON DELETE SET NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'business_profiles_verification_status_check'
    ) THEN
        ALTER TABLE business_profiles
            ADD CONSTRAINT business_profiles_verification_status_check
            CHECK (verification_status IN ('unsubmitted', 'pending', 'verified', 'rejected'));
    END IF;
END $$;

-- Backfill (only rows still at the default, so re-runs change nothing).
UPDATE business_profiles
SET verification_status = 'verified',
    verified_at = COALESCE(verified_at, updated_at)
WHERE is_verified = TRUE
  AND verification_status = 'unsubmitted';

UPDATE business_profiles bp
SET verification_status = 'pending',
    verification_submitted_at = COALESCE(bp.verification_submitted_at, u.updated_at)
FROM users u
WHERE u.id = bp.user_id
  AND u.status = 'active'
  AND bp.is_verified = FALSE
  AND bp.verification_status = 'unsubmitted';

CREATE INDEX IF NOT EXISTS idx_business_profiles_verification_queue
    ON business_profiles (verification_status, verification_submitted_at);

ALTER TABLE users ADD COLUMN IF NOT EXISTS previous_status user_status;

UPDATE campaigns SET status = 'paused', updated_at = NOW() WHERE status = 'pending';

COMMIT;

-- No admin account is created here. This file is applied to remote and
-- production databases, so it must never contain credentials.
-- Local development: run `make seed-admin` (migration/seed_admin.sql).
-- Other environments: create the first admin manually with a strong password.
