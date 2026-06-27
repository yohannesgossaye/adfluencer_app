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
